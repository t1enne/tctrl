package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ApiResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type UserConfig struct {
	User struct {
		Id string `json:"_id"`
	} `json:"user"`
	Token   string `json:"token"`
	ApiBase string `json:"apiBase"`
}

type BaseDto struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}
type UserHours struct {
	ID         string   `json:"_id"`
	CreatedAt  string   `json:"createdAt"`
	Date       string   `json:"date"`
	DeletedAt  string   `json:"deletedAt"`
	Hours      string   `json:"hours"`
	HoursTag   HoursTag `json:"hoursTag"`
	HoursTagID string   `json:"hoursTagId"`
	Notes      string   `json:"notes"`
	Release    Release  `json:"release"`
	ReleaseID  string   `json:"releaseId"`
	UpdatedAt  string   `json:"updatedAt"`
	UserID     string   `json:"userId"`
}
type HoursTag struct {
	ID           string `json:"_id"`
	CreatedAt    string `json:"createdAt"`
	DeletedAt    string `json:"deletedAt"`
	IconName     string `json:"iconName"`
	IsModifiable bool   `json:"isModifiable"`
	Name         string `json:"name"`
	UpdatedAt    string `json:"updatedAt"`
}
type Release struct {
	ID                     string   `json:"_id"`
	BillableHoursBudget    float64  `json:"billableHoursBudget"`
	BilledHours            *float64 `json:"billedHours"`
	ContingencyHoursBudget *float64 `json:"contingencyHoursBudget"`
	CreatedAt              string   `json:"createdAt"`
	DayCost                int      `json:"dayCost"`
	Deadline               string   `json:"deadline"`
	DeletedAt              string   `json:"deletedAt"`
	DeployedAt             string   `json:"deployedAt"`
	DocURL                 *string  `json:"docUrl"`
	HoursBudget            float64  `json:"hoursBudget"`
	IsArchived             bool     `json:"isArchived"`
	ManagementDeadline     string   `json:"managementDeadline"`
	ManagementDeployedAt   string   `json:"managementDeployedAt"`
	Name                   string   `json:"name"`
	Project                Project  `json:"project"`
	ProjectID              string   `json:"projectId"`
	Status                 string   `json:"status"`
	Typology               string   `json:"typology"`
	UpdatedAt              string   `json:"updatedAt"`
}
type Project struct {
	ID                  string   `json:"_id"`
	BillableHoursBudget *float64 `json:"billableHoursBudget"`
	CreatedAt           string   `json:"createdAt"`
	Customer            Customer `json:"customer"`
	CustomerDeadline    string   `json:"customerDeadline"`
	CustomerID          string   `json:"customerId"`
	Deadline            string   `json:"deadline"`
	DeletedAt           string   `json:"deletedAt"`
	HoursBudget         *float64 `json:"hoursBudget"`
	IsArchived          bool     `json:"isArchived"`
	Name                string   `json:"name"`
	ProjectManagerID    string   `json:"projectManagerId"`
	UpdatedAt           string   `json:"updatedAt"`
}
type Customer struct {
	ID        string `json:"_id"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	DeletedAt string `json:"deletedAt"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updatedAt"`
}
type Pagination struct {
	CurrentPage  int `json:"currentPage"`
	ItemsPerPage int `json:"itemsPerPage"`
	TotalItems   int `json:"totalItems"`
}

type AddHoursPayload struct {
	Notes      string `json:"notes"`
	Hours      string `json:"hours"`
	Date       string `json:"date"`
	ReleaseId  string `json:"releaseId"`
	HoursTagId string `json:"hoursTagId"`
	UserId     string `json:"userId"`
}

func readConfig(configPath string) []byte {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Panicf("Failed to read file %s\n", configPath)
	}
	return data
}

func parseConfig(data []byte) UserConfig {
	var c UserConfig
	err := json.Unmarshal(data, &c)
	if err != nil {
		log.Panicln("failed to parse config file")
	}
	return c
}

func GetConfig(p string) UserConfig {
	return parseConfig(readConfig(p))
}

func GetActiveTags(c UserConfig) []HoursTag {
	var data ApiResponse[HoursTag]
	Post("hoursTags/fb", "", c, &data)
	return Filter(data.Data, func(t HoursTag) bool { return t.Name != "Ferie e permessi" })
}

func GetCustomers(c UserConfig) []Customer {
	var d ApiResponse[Customer]
	Post("customers/fb", `{"order":{"name":"ASC"}}`, c, &d)
	return d.Data
}

func GetProjects(customerId string, c UserConfig) []Project {
	var d ApiResponse[Project]
	Post("projects/fb", `{"order":{"name":"ASC"}}`, c, &d)
	return d.Data
}

func GetReleases(projId string, c UserConfig) []Release {
	var data ApiResponse[Release]
	p := fmt.Sprintf(`{"order":{"name":"ASC"},"where":{"projectId":"%s"}}`, projId)
	Post("releases/fb", p, c, &data)
	return data.Data
}

func GetWorkedHours(payload string, c UserConfig) []UserHours {
	var data ApiResponse[UserHours]
	Post("userHours/fb", payload, c, &data)
	return data.Data
}

func AddHours(h AddHoursPayload, c UserConfig) {
	var r map[string]interface{}
	body, err := json.Marshal(h)
	if err != nil {
		log.Panicf("Failed to marshal hours: %s\n", err)
	}
	Post("userHours", string(body), c, &r)
}

func Post(url string, body string, c UserConfig, result interface{}) {
	request("POST", url, body, c, result)
}

func Delete(url string, body string, c UserConfig, result interface{}) {
	request("DELETE", url, body, c, result)
}

func request(method, url string, body string, c UserConfig, result interface{}) {
	fullUrl := c.ApiBase + url
	req, err := http.NewRequest(method, fullUrl, bytes.NewBufferString(body))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,it-IT;q=0.8,it;q=0.5,en;q=0.3")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Referer", "https://tpca.raintonic.com/")
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(c.Token))
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error performing request: %v\n", err)
		os.Exit(1)
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
}
