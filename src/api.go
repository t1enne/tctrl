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
type DayOff struct {
	ID         string   `json:"_id"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
	DeletedAt  *string  `json:"deletedAt"`
	UserID     string   `json:"userId"`
	HoursTagID string   `json:"hoursTagId"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
	Notes      string   `json:"notes"`
	Hours      string   `json:"hours"`
	Status     string   `json:"status"`
	User       User     `json:"user"`
	HoursTag   HoursTag `json:"hoursTag"`
}
type User struct {
	ID          string  `json:"_id"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	DeletedAt   *string `json:"deletedAt"`
	Email       string  `json:"email"`
	RoleID      string  `json:"roleId"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	IsDeletable bool    `json:"isDeletable"`
}
type AddDayOffPayload struct {
	User struct {
		ID string `json:"_id"`
	} `json:"user"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	Hours     float64 `json:"hours"`
	Notes     string  `json:"notes"`
	Status    string  `json:"status"`
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

func GetProjects(c UserConfig) []Project {
	var d ApiResponse[Project]
	Post("projects/fb", `{"order":{"name":"ASC"}}`, c, &d)
	return d.Data
}

func GetReleases(c UserConfig) []Release {
	var data ApiResponse[Release]
	p := `{"order":{"name":"ASC"},"include": ["project","project.customer"],"where": {"isArchived":false}}`
	Post("releases/fb", p, c, &data)
	return data.Data
}

func GetWorkedHours(payload string, c UserConfig) []UserHours {
	var data ApiResponse[UserHours]
	Post("userHours/fb", payload, c, &data)
	return data.Data
}

func GetDayOff(payload string, c UserConfig) []DayOff {
	var data ApiResponse[DayOff]
	Post("dayoffs/fb", payload, c, &data)
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

func AddDayOff(d AddDayOffPayload, c UserConfig) {
	var r map[string]interface{}
	body, err := json.Marshal(d)
	if err != nil {
		log.Panicf("Failed to marshal dayoff: %s\n", err)
	}
	Post("dayoffs", string(body), c, &r)
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
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Error performing request: %v\n", err)
		os.Exit(1)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		os.Exit(1)
	}
}
