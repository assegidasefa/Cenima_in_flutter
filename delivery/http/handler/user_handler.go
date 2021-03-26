package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joocosta/bloctrial/delivery/http/responses"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/permission"
	"github.com/joocosta/bloctrial/rtoken"
	"github.com/joocosta/bloctrial/user"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserHandler struct {
	service      user.UserService
	roleService  user.RoleService
	tokenService rtoken.Service
}

func NewUserHandler(us user.UserService, rs user.RoleService, ts rtoken.Service) *UserHandler {
	return &UserHandler{service: us, roleService: rs, tokenService:ts}
}

func (uh *UserHandler) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	// validate the token
	fn := func(w http.ResponseWriter, r *http.Request) {
		_token := r.Header.Get("Authorization")
		_token = strings.Replace(_token, "Bearer ", "", 1)
		valid, err := uh.tokenService.ValidateToken(_token)
		if err != nil && !valid {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthenticated: unauthorized to access the resource, log in again"))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)

}

func (uh *UserHandler) Authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_token := r.Header.Get("Authorization")
		_token = strings.Replace(_token, "Bearer ", "", 1)
		if len(_token) == 0 {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized: unauthorized to access the resource, log in again"))
			return
		}
		claim, err := uh.tokenService.GetClaims(_token)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		role, errs := uh.roleService.Role(claim.User.RoleID)
		if len(errs) > 0 || !permission.HasPermission(r.URL.Path, role.Name, r.Method) {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("not authorized to execute the request"))
			return
		}
		next.ServeHTTP(w, r)
		return
	}
}

func (uh *UserHandler) checkAdmin(roleID uint) bool {
	if roleID == 2 {
		return true
	}
	return false
}
func (uh *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if len(u.Email) == 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Email Should be not be empty"))
		return
	}
	if len(u.Password) == 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Password Should not be empty"))
		return
	}

	if uh.service.EmailExists(u.Email){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Email already occupied"))
		json.NewEncoder(w).Encode(err)
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Password Encryption  failed"))
		json.NewEncoder(w).Encode(err)
	}

	u.Password = string(pass)
	u.RoleID = 1

	newUser, errs := uh.service.StoreUser(&u)
	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("User Creation Failed"))
		return
	}

	json.NewEncoder(w).Encode(newUser)
}
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u model.User
	//var admin = false
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user, errs := uh.service.UserByEmail(u.Email)
	if len(errs) > 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("authentication error"))
		return
	}

	if uh.checkAdmin(user.RoleID) == true{
		user.Password = ""
		tokenString, err := uh.tokenService.GenerateToken(rtoken.CustomJwtClaim{
			User: *user,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().AddDate(0, 1, 1).Unix(),
				IssuedAt:  time.Now().Unix(),
				NotBefore: time.Now().Unix(),
			},
		})
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		})
		log.Println(user.Email + " has logged in!")
	} else{
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("authentication error"))
			return
		}

		//admin = uh.checkAdmin(user.RoleID)
		user.Password = ""
		tokenString, err := uh.tokenService.GenerateToken(rtoken.CustomJwtClaim{
			User: *user,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().AddDate(0, 1, 1).Unix(),
				IssuedAt:  time.Now().Unix(),
				NotBefore: time.Now().Unix(),
			},
		})
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		})
		log.Println(user.Email + " has logged in!")
	}
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_token := r.Header.Get("Authorization")
	_token = strings.Replace(_token, "Bearer ", "", 1)
	claim, err := uh.tokenService.GetClaims(_token)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	log.Println(time.Now().Unix() - claim.ExpiresAt)
	claim.ExpiresAt = time.Now().Unix()
	log.Println(claim.ExpiresAt)
	log.Println("User logged out")
}
