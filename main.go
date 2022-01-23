package main

import (
	"strings"
	"strconv"
	"reflect"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)


type pokemon struct {
	name		 string
	description	 string
	habitat		 string
	is_legendary	 bool
}

func main(){
 	router := gin.Default()
	router.GET("/pokemon/:name", getPokemon)
	router.GET("/home/")
	router.Run("localhost:8080")
}

func getRequest(url string) string {
	client := &http.Client{Timeout: time.Second * 2}
	
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}
	
	req.Header.Set("User-Agent", "pokedex")

	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	t := reflect.TypeOf(resBody)
	k := t.Kind()
	log.Info("type readall: ", t, "kind: ", k) 
	returnResponse := string(resBody)
	return returnResponse
}




func getHabitat(pokemonName string) string{
	url := "https://pokeapi.co/api/v2/pokemon-habitat/"
	// from using api, turns out there's 9 habitats
	for i := 0; i<10; i++ {
		habitat := getRequest(url+strconv.Itoa(i))
		// if the pokemon is in it, the name will be in the habitat string
		if strings.Contains(habitat, pokemonName){
			resBytes := []byte(habitat)

			var jsonRes map[string]interface{}
			err := json.Unmarshal(resBytes, &jsonRes)
			if err != nil {
				log.Info("could not parse response")
				log.Fatal(err)
			}			
			return jsonRes["name"].(string) 
		}
	}
	return ""
}

func parseToJson(response string) map[string]interface{} {

	resBytes := []byte(response)

	var jsonRes map[string]interface{}
	err := json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Info("could not parse response")
		log.Fatal(err)
	}
	return jsonRes			
}

func unpackDescription(flavorEntries map[string]interface{}) string {
	log.Info("flavor entries: ", flavorEntries["flavor_text_entries"])
	// TODO: UNPACK THE DESCRIPTION -> IT'S A LIST OF DICTIONARIES, WE JUST WANT ONE KEY OF THE FIRST ELEMENT
	/*for f := range flavorEntries["flavor_text_entries"] {
		log.Info("flavor", f)
	
	}*/
	return ""
	/*
	var flavours map[string]interface{}
	_ = json.Unmarshal([]byte(flavorEntries), &flavours)
	log.Info("unmarshaled: %v", flavours)
	*/
}

func getIsLegendaryAndDescription(pokemonName string) (bool, string){
	url := "https://pokeapi.co/api/v2/pokemon-species/" + pokemonName
	response := getRequest(url)
 	jsonResponse := parseToJson(response)

	is_legendary := jsonResponse["is_legendary"].(bool)	

	//description := jsonResponse["flavor_text_entries"][0]["flavor"].(string)
	unpackDescription(jsonResponse)
	description := ""
	return is_legendary, description
}

/*
func formatAsPokemon(jsonResp map[string]interface{}) pokemon{
	name := jsonResp["name"].(string)
	description := "description"
	habitat := getHabitat(name)
	is_legendary := "description"
	//var p pokemon{jsonResp}
}
*/

func getPokemon(c *gin.Context){

	log.Info("getting pokemon!")
	url := "http://pokeapi.co/api/v2/pokemon/"
	name := c.Param("name")
	url += name
	log.Print("Fetching from url " + url)	

	response := getRequest(url)
	//log.Info("raw response: " + response)		
	jsonRes := parseResponse(response)
	log.Info(jsonRes)

	h := getHabitat(name)
	l, d := getIsLegendaryAndDescription(name)

	log.Info("habitat: ", h)
	log.Info("is_legendary", l)	
	log.Info("description", d)
	c.IndentedJSON(http.StatusOK, gin.H{"hello":"world"})
}

func parseResponse(response string) map[string]interface{}{
	resBytes := []byte(response)
	var jsonRes map[string]interface{}
	err := json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Info("could not parse response")
		log.Fatal(err)
	}
	return jsonRes

}


