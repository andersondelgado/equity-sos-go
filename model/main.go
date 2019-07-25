package model

// import "database/sql"

import (
	"time"
)

/////////begin:get document by fields///////////////////////
type QuerySelectorPaginateIndex struct {
	Selector interface{} `json:"selector"`
	Limit    uint        `json:"limit"`
	Skip     uint        `json:"skip"`
	UseIndex string      `json:"use_index"`
	Fields   []string    `json:"fields"`
}

type QuerySelectorPaginate struct {
	Selector interface{} `json:"selector"`
	Limit    uint        `json:"limit"`
	Skip     uint        `json:"skip"`
	Fields   []string    `json:"fields"`
}

type QuerySelectorAll struct {
	Selector interface{} `json:"selector"`
	Fields   []string    `json:"fields"`
}

/////////end:get document by fields///////////////////////

/////////begin:all doc result///////////////////////
type AlldocsResult struct {
	TotalRows int `json:"total_rows"`
	Offset    int
	Rows      []map[string]interface{}
}

/////////end:all doc result///////////////////////

/////////begin:test document///////////////////////
type TestDocument struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc TestDoc `json:"doc"`
}

type TestDocumentsArray struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc      []TestDoc `json:"docs"`
	Bookmark string    `json:"bookmark"`
	Warning  string    `json:"warning"`
}

type TestDoc struct {
	ID   string `json:"_id"`
	Rev  string `json:"_rev"`
	Test Test   `json:"tests"`
}

type Test struct {
	// ID          uint      `json:"id"`
	IDs         string    `json:"_id"`
	Rev         string    `json:"_rev"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

/////////end:test document///////////////////////

/////////begin:user document///////////////////////
type UserDocument struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc UserDoc `json:"doc"`
}

type UserDocumentsArray struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc      []UserDoc `json:"docs"`
	Bookmark string    `json:"bookmark"`
	Warning  string    `json:"warning"`
}

type UserDoc struct {
	ID   string `json:"_id"`
	Rev  string `json:"_rev"`
	User User   `json:"users"`
}

type User struct {
	// ID        uint      `json:"id"`
	IDs       string    `json:"_id"`
	Rev       string    `json:"_rev"`
	ID        string    `json:"id"`
	Avatar    string    `json:"avatar"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PayloadUserRoles struct {
	Rol      string `json:"rol"`
	Avatar   string `json:"avatar"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

/////////end:document///////////////////////

/////////begin:permission document///////////////////////
type PermissionDocument struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc PermissionDoc `json:"doc"`
}

type PermissionDocumentsArray struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc      []PermissionDoc `json:"docs"`
	Bookmark string          `json:"bookmark"`
	Warning  string          `json:"warning"`
}

type PermissionDoc struct {
	ID         string     `json:"_id"`
	Rev        string     `json:"_rev"`
	Permission Permission `json:"permissions"`
}

type Permission struct {
	IDs                    string      `json:"_id"`
	Rev                    string      `json:"_rev"`
	ID                     string      `json:"id"`
	ObjectModulePermission string      `json:"object_module_permission"`
	Rol                    string      `json:"rol"`
	Modules                interface{} `json:"module_permission"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
}

/////////end:permission document///////////////////////

/////////begin:profile document///////////////////////
type ProfileDocument struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc ProfileDoc `json:"doc"`
}

type ProfileDocumentsArray struct {
	// ID  string  `json:"_id"`
	// Rev string  `json:"_rev"`
	Doc      []ProfileDoc `json:"docs"`
	Bookmark string       `json:"bookmark"`
	Warning  string       `json:"warning"`
}

type ProfileDoc struct {
	ID      string  `json:"_id"`
	Rev     string  `json:"_rev"`
	Profile Profile `json:"profiles"`
}

type Profile struct {
	IDs          string     `json:"_id"`
	Rev          string     `json:"_rev"`
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	PermissionID string     `json:"permission_id"`
	User         User       `json:"user"`
	Permission   Permission `json:"permission"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

/////////end:profile document///////////////////////

/////////begin:status document///////////////////////
type StatusDocument struct {
	Doc StatusDoc `json:"doc"`
}
type StatusDocumentsArray struct {
	Doc      []StatusDoc `json:"docs"`
	Bookmark string      `json:"bookmark"`
	Warning  string      `json:"warning"`
}
type StatusDoc struct {
	ID     string `json:"_id"`
	Rev    string `json:"_rev"`
	Status Status `json:"status"`
}
type Status struct {
	IDs         string    `json:"_id"`
	Rev         string    `json:"_rev"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

/////////end:status document///////////////////////

/////////begin:category document///////////////////////
type CategoryDocument struct {
	Doc CategoryDoc `json:"doc"`
}
type CategoryDocumentsArray struct {
	Doc      []CategoryDoc `json:"docs"`
	Bookmark string        `json:"bookmark"`
	Warning  string        `json:"warning"`
}
type CategoryDoc struct {
	ID       string   `json:"_id"`
	Rev      string   `json:"_rev"`
	Category Category `json:"categories"`
}
type Category struct {
	IDs       string    `json:"_id"`
	Rev       string    `json:"_rev"`
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/////////end:category document///////////////////////

/////////begin:KYC document///////////////////////
type KYCDocument struct {
	Doc KYCDoc `json:"doc"`
}
type KYCDocumentsArray struct {
	Doc      []KYCDoc `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}
type KYCDoc struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
	KYC KYC    `json:"kyc"`
}
type KYC struct {
	IDs         string    `json:"_id"`
	Rev         string    `json:"_rev"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

/////////end:KYC document///////////////////////

/////////begin:user KYC document///////////////////////
type UserKYCDocument struct {
	Doc UserKYCDoc `json:"doc"`
}
type UserKYCDocumentsArray struct {
	Doc []UserKYCDoc `json:"docs"`

	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}
type UserKYCDoc struct {
	ID      string  `json:"_id"`
	Rev     string  `json:"_rev"`
	UserKYC UserKYC `json:"user_kyc"`
}
type UserKYC struct {
	IDs       string    `json:"_id"`
	Rev       string    `json:"_rev"`
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	KYCID     string    `json:"kyc_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/////////end:user KYC document///////////////////////

/////////begin:country document///////////////////////
type CountrysDocument struct {
	Doc CountrysDoc `json:"doc"`
}
type CountrysDocumentsArray struct {
	Doc      []CountrysDoc `json:"docs"`
	Bookmark string        `json:"bookmark"`
	Warning  string        `json:"warning"`
}

type CountrysDoc struct {
	ID       string   `json:"_id"`
	Rev      string   `json:"_rev"`
	Countrys Countrys `json:"countrys"`
}
type Countrys struct {
	IDs       string    `json:"_id"`
	Rev       string    `json:"_rev"`
	ID        string    `json:"id"`
	Short     string    `json:"short"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/////////end:country document///////////////////////

/////////begin:article document///////////////////////
type ArticleDocument struct {
	Doc ArticleDoc `json:"doc"`
}

type ArticleDocumentsArray struct {
	Doc      []ArticleDoc `json:"docs"`
	Bookmark string       `json:"bookmark"`
	Warning  string       `json:"warning"`
}
type ArticleDoc struct {
	ID      string  `json:"_id"`
	Rev     string  `json:"_rev"`
	Article Article `json:"articles"`
}
type Article struct {
	IDs          string    `json:"_id"`
	Rev          string    `json:"_rev"`
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TypeArticle  string    `json:"type_article"`
	Composition  string    `json:"composition"`
	Presentation string    `json:"presentation"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

/////////end:article document///////////////////////

/////////begin:post document///////////////////////
type PostsDocument struct {
	Doc PostsDoc `json:"doc"`
}
type PostsDocumentsArray struct {
	Doc      []PostsDoc `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}
type PostsDoc struct {
	ID    string `json:"_id"`
	Rev   string `json:"_rev"`
	Posts Posts  `json:"posts"`
}
type Posts struct {
	IDs          string         `json:"_id"`
	Rev          string         `json:"_rev"`
	ID           string         `json:"id"`
	PostsData    PostsData      `json:"post_data"`
	ArticlesData []ArticlesData `json:"articles_data"`
	StatusPost   []StatusPost   `json:"status_post"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
type PostsData struct {
	UserID         string           `json:"user_id"`
	RangeDamage    string           `json:"range_damage"`
	HashSolidty    string           `json:"hash_solidty"`
	FullText       string           `json:"full_text"`
	LatLong        LatLong          `json:"lat_long"`
	Atachments     []Atachments     `json:"atachments"`
	CategoryPostID []string         `json:"category_post_id"`
	Category       []Category       `json:"category_post"`
	PromedyPeoples []PromedyPeoples `json:"promedy_peoples"`
}
type LatLong struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PromedyPeoples struct {
	People   string `json:"people"`
	Quantity uint   `json:"quantity"`
}
type ArticlesData struct {
	ArticleID      string           `json:"article_id"`
	Name           string           `json:"name"`
	QuantityAsk    uint             `json:"quantity_ask"`
	QuantityLeft   uint             `json:"quantity_left"`
	ArticlesActors []ArticlesActors `json:"articles_actors"`
}
type ArticlesActors struct {
	SequencesID            string    `json:"sequences_id"`
	UserID                 string    `json:"user_id"`
	User                   User      `json:"users"`
	QuantityDelivery       uint      `json:"quantity_delivery"`
	QuantityLeft           uint      `json:"quantity_left"`
	StatusDeliverySender   bool      `json:"status_delivery_sender"`
	StatusDeliveryReciever bool      `json:"status_delivery_reciever"`
	StatusOrder            string    `json:"status_order"`
	CreatedAt              time.Time `json:"created_at"`
	DateDeliverySender     time.Time `json:"date_delivery_sender"`
	DateDeliveryReciever   time.Time `json:"date_delivery_reciever"`
}
type StatusPost struct {
	StatusID  string    `json:"status_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Atachments struct {
	Name []string `json:"name"`
	Type string   `json:"type"`
}

type PostChaincode struct {
}
/////////end:post document///////////////////////

/////////begin:document///////////////////////
type Document struct {
	Doc interface{} `json:"doc"`
}

/////////end:document///////////////////////

/////////begin:module///////////////////////
type Module struct {
	Url          string   `json:"url"`
	IconWeb      string   `json:"iconWeb"`
	Mobile       string   `json:"mobile"`
	Android      string   `json:"android"`
	Ios          string   `json:"ios"`
	Ionic        string   `json:"ionic"`
	Name         string   `json:"name"`
	LangProperty string   `json:"lang_property"`
	Acl          string   `json:"acl"`
	Visible      bool     `json:"visible"`
	Value        []string `json:"value"`
}

type ObjectModulePermission struct {
	IDs    string   `json:"_id"`
	Rev    string   `json:"_rev"`
	Rol    string   `json:"rol"`
	Module []Module `json:"module"`
}

/////////end:module///////////////////////

/////////begin:roles dev///////////////////////
type RolesDev struct {
	IDs    string      `json:"_id"`
	Rev    string      `json:"_rev"`
	Rol    string      `json:"rol"`
	Module interface{} `json:"module"`
}

/////////end:roles dev///////////////////////

/////////begin:d///////////////////////
/////////end:d///////////////////////
