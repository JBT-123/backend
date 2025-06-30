package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/How-to-get-ABG/backend/internal/models"
)

///URL Might not need to be in there


func UpdateUserInfo(w http.ResponseWriter, r *http.Request){
	gender_enums:=[]string{"male", "female", "nonbinary", "other"}
	var req models.UserInfo
	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w ,"Invalid request body:", http.StatusBadRequest)
	}

	idVal := r.Context().Value("user_id")

	id, ok := idVal.(int)
	if !ok {
		fmt.Println("Invalid or missing user_id in context")
		return;
	}	
	
	if gender:= strings.ToLower(req.Gender); gender != "" && !slices.Contains(gender_enums, gender){
		http.Error(w, "Invalid Gender Input:", http.StatusBadRequest)
		return;
	}
	// add a check for interests in the future when updating preferences

	if err:=db.UpdateUserDB(&req,id); err!=nil{
		http.Error(w,"Error in updating user Info", http.StatusBadRequest)
		return;
	}
	
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)

json.NewEncoder(w).Encode(map[string]interface{}{
    "msg":     "User updated successfully",
    "success": true,
})
}

func UpdateUserPref(w http.ResponseWriter, r *http.Request){
	gender_enums:=[]string{"male", "female", "nonbinary", "other"}
	var req models.UserPref
	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w ,"Invalid request body:", http.StatusBadRequest)
	}

	idVal := r.Context().Value("user_id")

	id, ok := idVal.(int)
	if !ok {
		fmt.Println("Invalid or missing user_id in context")
		return;
	}	
	
	if gender:= strings.ToLower(req.PrefGender); gender != "" && !slices.Contains(gender_enums, gender){
		http.Error(w, "Invalid Gender Input:", http.StatusBadRequest)
		return;
	}
	// add a check for interests in the future when updating preferences

	if err:=db.UpdateUserPref(&req,id); err!=nil{
		http.Error(w,"Error in updating user Info", http.StatusBadRequest)
		return;
	}
	
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)

json.NewEncoder(w).Encode(map[string]interface{}{
    "msg":     "User Pref updated successfully",
    "success": true,
})
}

func Connect(w http.ResponseWriter, r *http.Request){
	idVal := r.Context().Value("user_id")

	id, ok := idVal.(int)
	if !ok {
		fmt.Println("Invalid or missing user_id in context")
		return;
	}	
	connectData, err:= db.UserLinkedPref(id);
	if err!=nil{
		http.Error(w, "An error has occurred: ", http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(connectData)
}

