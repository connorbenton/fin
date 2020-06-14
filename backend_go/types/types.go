package types

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name" db:"name"`
	Institution        string          `json:"institution" db:"institution"`
	IgnoreTransactions bool            `json:"ignore_transactions" db:"ignore_transactions"`
	AccountID          string          `json:"account_id" db:"account_id"`
	ItemID             string          `json:"item_id" db:"item_id"`
	Type               string          `json:"type" db:"type"`
	Subtype            string          `json:"subtype" db:"subtype"`
	Balance            decimal.Decimal `json:"balance" db:"balance"`
	Limit              decimal.Decimal `json:"limit" db:"limit"`
	Available          decimal.Decimal `json:"available" db:"available"`
	Currency           string          `json:"currency" db:"currency"`
	Provider           string          `json:"provider" db:"provider"`
	RunningTotal       decimal.Decimal `json:"running_total" db:"running_total"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
}

type ItemToken struct {
	ID          int    `json:"id"`
	Institution string `json:"institution" db:"institution"`
	AccessToken string `json:"-" db:"access_token"`
	ItemID      string `json:"item_id" db:"item_id"`
	Provider    string `json:"provider" db:"provider"`
	Interactive                bool      `json:"interactive" db:"interactive"`
	NeedsReLogin               bool      `json:"needs_re_login" db:"needs_re_login"`
	LastRefresh                time.Time `json:"last_refresh" db:"last_refresh"`
	NextRefreshPossible        time.Time `json:"next_refresh_possible" db:"next_refresh_possible"`
	LastDownloadedTransactions time.Time `json:"last_downloaded_transactions" db:"last_downloaded_transactions"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at" db:"updated_at"`
}

type Category struct {
	ID                  int       `json:"id"`
	TopCategory         string    `json:"top_category" db:"top_category"`
	SubCategory         string    `json:"sub_category" db:"sub_category"`
	ExcludeFromAnalysis bool      `json:"exclude_from_analysis" db:"exclude_from_analysis"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	Count               int
	Total               decimal.Decimal
}

type CategorySE struct {
	ID             int       `json:"id"`
	TopCategory    string    `json:"top_category" db:"top_category"`
	SubCategory    string    `json:"sub_category" db:"sub_category"`
	BottomCategory string    `json:"bottom_category" db:"bottom_category"`
	LinkToAppCat   int       `json:"link_to_app_cat" db:"link_to_app_cat"`
	AppCatName     string    `json:"app_cat_name" db:"app_cat_name"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CategoryPlaid struct {
	ID           int       `json:"id"`
	Hierarchy    string    `json:"hierarchy" db:"hierarchy"`
	CatID        string    `json:"cat_i_d" db:"cat_i_d"`
	LinkToAppCat int       `json:"link_to_app_cat" db:"link_to_app_cat"`
	AppCatName   string    `json:"app_cat_name" db:"app_cat_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID                  int             `json:"id"`
	Date                string          `json:"date" db:"date"`
	TransactionID       string          `json:"transaction_id" db:"transaction_id"`
	Description         string          `json:"description" db:"description"`
	OriginalDescription string          `json:"original_description" db:"original_description"`
	Amount              decimal.Decimal `json:"amount" db:"amount"`
	NormalizedAmount    decimal.Decimal `json:"normalized_amount" db:"normalized_amount"`
	TransactionType     string          `json:"transaction_type" db:"transaction_type"`
	Category            int             `json:"category" db:"category"`
	CategoryName        string          `json:"category_name" db:"category_name"`
	AccountName         string          `json:"account_name" db:"account_name"`
	CurrencyCode        string          `json:"currency_code" db:"currency_code"`
	AccountID           string          `json:"account_id" db:"account_id"`
	Labels              string          `json:"labels" db:"labels"`
	Notes               string          `json:"notes" db:"notes"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
}

var TreeRanges = [8]string{
	"last30",
	"thisMonth",
	"lastMonth",
	"last6Months",
	"thisYear",
	"lastYear",
	"fromBeginning",
	"custom",
}

type Tree struct {
	Name      string `json:"name" db:"name"`
	FirstDate string `json:"first_date" db:"first_date"`
	LastDate  string `json:"last_date" db:"last_date"`
	Data         string    `json:"data" db:"data"`
	DataNoInvest string    `json:"data_no_invest" db:"data_no_invest"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type TreeData struct {
	Name     string          `json:"name"`
	Children []ChildTop      `json:"children"`
	Value    decimal.Decimal `json:"value"`
	Count    int64           `json:"count"`
	TrueCount   int64           `json:"trueCount"`
	Per30       decimal.Decimal `json:"per30"`
	IncomeTotal decimal.Decimal `json:"income_total"`
}

type ChildTop struct {
	Name     string          `json:"name"`
	Children []ChildSub      `json:"children"`
	Value    decimal.Decimal `json:"value"`
	Count    int64           `json:"count"`
	TrueCount int64           `json:"trueCount"`
	DbID      int             `json:"dbID"`
	Percent   string          `json:"percent"`
	Per30     decimal.Decimal `json:"per30"`
}

type ChildSub struct {
	Name    string          `json:"name"`
	DbID    int             `json:"dbID"`
	Value   decimal.Decimal `json:"value"`
	Count   int64           `json:"count"`
	Percent string          `json:"percent"`
	TrueCount int64           `json:"trueCount"`
	Per30     decimal.Decimal `json:"per30"`
}

type CustomRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

var MintCatMap = map[string]int{
	"Buy":                        69,
	"Investments":                99,
	"Dividend & Cap Gains":       63,
	"Sports":                     51,
	"Coffee Shops":               37,
	"Gift":                       42,
	"Shipping":                   76,
	"Hide from Budgets & Trends": 109,
	"Business Services":          76,
	"Home Improvement":           54,
	"Home Services":              54,
	"Withdrawal":                 107,
	"Office Supplies":            57,
	"Pets":                       76,
	"Credit Card Payment":        100,
	"Brokerage":                  99,
	"Interest Income":            63,
	"Brokerage Investment":       99,
	"Printing":                   76,
}

type MatchingAccount struct {
	ImportKey      string `json:"importKey"`
	RefAccountID   string `json:"refAccountID"`
	RefAccountName string `json:"refAccountName"`
}

type CompareTrans struct {
	Date         string          `json:"date"`
	Description  string          `json:"description"`
	Amount       decimal.Decimal `json:"amount"`
	AccountName  string          `json:"account_name"`
	AccountID    string          `json:"account_id"`
	CurrencyCode string          `json:"currency_code"`
}

type CompareTransSingle struct {
	Trans1  CompareTrans `json:"trans1"`
	Trans2  CompareTrans `json:"trans2"`
	IsMatch bool         `json:"isMatch`
	Type    string       `json:"type"`
}

type CompareCatsSingle struct {
	Category        string `json:"category"`
	AssignedCat     int    `json:"assignedCat"`
	AssignedCatName string `json:"assignedCatName"`
}

type ImportTransaction struct {
	Date                string          `json:"date"`
	Description         string          `json:"description"`
	OriginalDescription string          `json:"originalDescription"`
	Amount              decimal.Decimal `json:"amount"`
	TransactionType     string          `json:"transactionType"`
	Category            string          `json:"category"`
	AccountName         string          `json:"accountName"`
	Labels              string          `json:"labels"`
	Notes               string          `json:"notes"`
	CurrencyCode        string          `json:"currency_code"`
}

type ImportPostData struct {
	Catres             []CompareCatsSingle `json:"catres"`
	IdentifiedAccounts []MatchingAccount   `json:"identifiedAccounts"`
	TxSet              []ImportTransaction `json:"transactions"`
}

type GenerateTokenPost struct {
	ItemID string `json:"item_id"`
}

type CreateTokenPost struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

type EcbFXRate struct {
	Date struct {
		Value string `xml:"value,attr"`
	} `xml:"ObsDimension"`
	Rate struct {
		Value string `xml:"value,attr"`
	} `xml:"ObsValue"`
}

type EcbFXCurrency struct {
	SeriesKey struct {
		Text   string `xml:",chardata"`
		Values []struct {
			ID    string `xml:"id,attr"`
			Value string `xml:"value,attr"`
		} `xml:"Value"`
	} `xml:"SeriesKey"`
	Rates []EcbFXRate `xml:"Obs"`
}

type EcbFX struct {
	Currencies []EcbFXCurrency `xml:"DataSet>Series"`
}

// Fx type exported for currency rate lookups
type Fx struct {
	FxDate time.Time       `db:"fx_date"`
	Rate   decimal.Decimal `db:"rate"`
}

type SEConnection struct {
	ID                      string    `json:"id"`
	Secret                  string    `json:"secret"`
	ProviderID              string    `json:"provider_id"`
	ProviderCode            string    `json:"provider_code"`
	ProviderName            string    `json:"provider_name"`
	CustomerID              string    `json:"customer_id"`
	NextRefreshPossibleAt   time.Time `json:"next_refresh_possible_at"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	Status                  string    `json:"status"`
	Categorization          string    `json:"categorization"`
	DailyRefresh            bool      `json:"daily_refresh"`
	StoreCredentials        bool      `json:"store_credentials"`
	CountryCode             string    `json:"country_code"`
	LastSuccessAt           time.Time `json:"last_success_at"`
	ShowConsentConfirmation bool      `json:"show_consent_confirmation"`
	LastConsentID           string    `json:"last_consent_id"`
	LastAttempt             struct {
		ID                   string      `json:"id"`
		Finished             bool        `json:"finished"`
		APIMode              string      `json:"api_mode"`
		APIVersion           string      `json:"api_version"`
		Locale               string      `json:"locale"`
		UserPresent          bool        `json:"user_present"`
		CustomerLastLoggedAt interface{} `json:"customer_last_logged_at"`
		RemoteIP             string      `json:"remote_ip"`
		FinishedRecent       bool        `json:"finished_recent"`
		Partial              bool        `json:"partial"`
		AutomaticFetch       bool        `json:"automatic_fetch"`
		DailyRefresh         bool        `json:"daily_refresh"`
		Categorize           bool        `json:"categorize"`
		CustomFields         struct {
		} `json:"custom_fields"`
		DeviceType              string        `json:"device_type"`
		UserAgent               string        `json:"user_agent"`
		ExcludeAccounts         []interface{} `json:"exclude_accounts"`
		FetchScopes             []string      `json:"fetch_scopes"`
		FromDate                string        `json:"from_date"`
		ToDate                  string        `json:"to_date"`
		Interactive             bool          `json:"interactive"`
		StoreCredentials        bool          `json:"store_credentials"`
		IncludeNatures          interface{}   `json:"include_natures"`
		ShowConsentConfirmation bool          `json:"show_consent_confirmation"`
		ConsentID               string        `json:"consent_id"`
		FailAt                  interface{}   `json:"fail_at"`
		FailMessage             interface{}   `json:"fail_message"`
		FailErrorClass          interface{}   `json:"fail_error_class"`
		CreatedAt               time.Time     `json:"created_at"`
		UpdatedAt               time.Time     `json:"updated_at"`
		SuccessAt               time.Time     `json:"success_at"`
		LastStage               struct {
			ID                       string      `json:"id"`
			Name                     string      `json:"name"`
			InteractiveHTML          interface{} `json:"interactive_html"`
			InteractiveFieldsNames   interface{} `json:"interactive_fields_names"`
			InteractiveFieldsOptions interface{} `json:"interactive_fields_options"`
			CreatedAt                time.Time   `json:"created_at"`
			UpdatedAt                time.Time   `json:"updated_at"`
		} `json:"last_stage"`
	} `json:"last_attempt"`
}

type SEAccount struct {
	ID           string          `json:"id"`
	ConnectionID string          `json:"connection_id"`
	Name         string          `json:"name"`
	Nature       string          `json:"nature"`
	Balance      decimal.Decimal `json:"balance"`
	CurrencyCode string          `json:"currency_code"`
	Extra        struct {
		Iban              string          `json:"iban"`
		ClientName        string          `json:"client_name"`
		AccountName       string          `json:"account_name"`
		CardType          string          `json:"card_type"`
		CreditLimit       decimal.Decimal `json:"credit_limit"`
		AvailableAmount   decimal.Decimal `json:"available_amount"`
		TransactionsCount struct {
			Posted  int `json:"posted"`
			Pending int `json:"pending"`
		} `json:"transactions_count"`
		LastPostedTransactionID string `json:"last_posted_transaction_id"`
	} `json:"extra,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SETransaction struct {
	ID           string          `json:"id"`
	AccountID    string          `json:"account_id"`
	Duplicated   bool            `json:"duplicated"`
	Mode         string          `json:"mode"`
	Status       string          `json:"status"`
	MadeOn       string          `json:"made_on"`
	Amount       decimal.Decimal `json:"amount"`
	CurrencyCode string          `json:"currency_code"`
	Description  string          `json:"description"`
	Category     string          `json:"category"`
	Extra        struct {
		PostingDate              string          `json:"posting_date"`
		ClosingBalance           decimal.Decimal `json:"closing_balance"`
		AccountBalanceSnapshot   decimal.Decimal `json:"account_balance_snapshot"`
		CategorizationConfidence float64         `json:"categorization_confidence"`
	} `json:"extra,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountResponse struct {
	Data []SEAccount `json:"data"`
	Meta struct {
		NextID   interface{} `json:"next_id"`
		NextPage interface{} `json:"next_page"`
	} `json:"meta"`
}

type ConnectionResponse struct {
	Data []SEConnection `json:"data"`
	Meta struct {
		NextID   interface{} `json:"next_id"`
		NextPage interface{} `json:"next_page"`
	} `json:"meta"`
}

type TransactionsResponse struct {
	Data []SETransaction `json:"data"`
	Meta struct {
		NextID   interface{} `json:"next_id"`
		NextPage interface{} `json:"next_page"`
	} `json:"meta"`
}

type CreateRefreshResponse struct {
	Data struct {
		ExpiresAt  time.Time `json:"expires_at"`
		ConnectURL string    `json:"connect_url"`
	} `json:"data"`
}

func PrepTransSt(txn *sqlx.Tx) *sqlx.NamedStmt {
	tquery := `INSERT INTO transactions('date', transaction_id, description, amount, normalized_amount, category,
				category_name, account_name, currency_code, account_id)
				VALUES(:date, :transaction_id, :description, :amount, :normalized_amount, :category,
				:category_name, :account_name, :currency_code, :account_id) 
				ON CONFLICT (transaction_id) DO UPDATE SET
				'date' = excluded.'date',
				description = excluded.description,
				amount = excluded.amount,
				normalized_amount = excluded.normalized_amount`
	tstmt, err := txn.PrepareNamed(tquery)
	if err != nil {
		panic(err)
	}
	return tstmt
}

func PrepTransUpsertSt(txn *sqlx.Tx) *sqlx.NamedStmt {
	tquery := `INSERT INTO transactions('date', transaction_id, description, amount, normalized_amount, category,
				category_name, account_name, currency_code, account_id)
				VALUES(:date, :transaction_id, :description, :amount, :normalized_amount, :category,
				:category_name, :account_name, :currency_code, :account_id) 
				ON CONFLICT (transaction_id) DO UPDATE SET
				'date' = excluded.'date',
				description = excluded.description,
				amount = excluded.amount,
				normalized_amount = excluded.normalized_amount,
				category = excluded.category,
				category_name = excluded.category_name`
	tstmt, err := txn.PrepareNamed(tquery)
	if err != nil {
		panic(err)
	}
	return tstmt
}

func PrepAccountSt(txn *sqlx.Tx) *sqlx.NamedStmt {
	aquery := `INSERT INTO accounts(name, institution, provider, account_id, item_id, type, 'limit', available, balance, currency, subtype)
				VALUES(:name, :institution, :provider, :account_id, :item_id, :type, :limit, :available, :balance, :currency, :subtype) 
				ON CONFLICT (account_id, provider) DO UPDATE SET
				'limit' = excluded.'limit',
				available = excluded.available,
				balance = excluded.balance`
	astmt, err := txn.PrepareNamed(aquery)
	if err != nil {
		panic(err)
	}
	return astmt
}

func PrepAccountUpsertIgnoreSt(txn *sqlx.Tx) *sqlx.NamedStmt {
	aquery := `INSERT INTO accounts(name, institution, provider, account_id, item_id, type, 'limit', available, balance, currency, subtype, ignore_transactions)
				VALUES(:name, :institution, :provider, :account_id, :item_id, :type, :limit, :available, :balance, :currency, :subtype, :ignore_transactions) 
				ON CONFLICT (account_id, provider) DO UPDATE SET
				ignore_transactions = excluded.ignore_transactions`
	astmt, err := txn.PrepareNamed(aquery)
	if err != nil {
		panic(err)
	}
	return astmt
}

func PrepAccountUpsertNameSt(txn *sqlx.Tx) *sqlx.NamedStmt {
	aquery := `INSERT INTO accounts(name, institution, provider, account_id, item_id, type, 'limit', available, balance, currency, subtype, ignore_transactions)
				VALUES(:name, :institution, :provider, :account_id, :item_id, :type, :limit, :available, :balance, :currency, :subtype, :ignore_transactions) 
				ON CONFLICT (account_id, provider) DO UPDATE SET
				name = excluded.name`
	astmt, err := txn.PrepareNamed(aquery)
	if err != nil {
		panic(err)
	}
	return astmt
}

func PrepItemSt(txn *sqlx.Tx) *sqlx.NamedStmt {
	iquery := `INSERT INTO item_tokens(institution, provider, interactive, last_refresh, next_refresh_possible, item_id, needs_re_login, access_token, last_downloaded_transactions)
				VALUES(:institution, :provider, :interactive, :last_refresh, :next_refresh_possible, :item_id, :needs_re_login, :access_token, :last_downloaded_transactions) 
				ON CONFLICT (item_id, provider) DO UPDATE SET
				interactive = excluded.interactive,
				last_refresh = excluded.last_refresh,
				next_refresh_possible = excluded.next_refresh_possible,
				needs_re_login = excluded.needs_re_login,
				last_downloaded_transactions = excluded.last_downloaded_transactions`
	istmt, err := txn.PrepareNamed(iquery)
	if err != nil {
		panic(err)
	}
	return istmt
}

func PrepItemStOnlyTx(txn *sqlx.Tx) *sqlx.NamedStmt {
	iquery := `INSERT INTO item_tokens(institution, provider, interactive, last_refresh, next_refresh_possible, item_id, needs_re_login, access_token, last_downloaded_transactions)
				VALUES(:institution, :provider, :interactive, :last_refresh, :next_refresh_possible, :item_id, :needs_re_login, :access_token, :last_downloaded_transactions) 
				ON CONFLICT (item_id, provider) DO UPDATE SET
				last_downloaded_transactions = excluded.last_downloaded_transactions`
	istmt, err := txn.PrepareNamed(iquery)
	if err != nil {
		panic(err)
	}
	return istmt
}

func PrepTreeSt(txn *sqlx.Tx) *sqlx.NamedStmt {
	iquery := `INSERT INTO analysis_trees(name, first_date, last_date, data, data_no_invest)
				VALUES(:name, :first_date, :last_date, :data, :data_no_invest) 
				ON CONFLICT (name) DO UPDATE SET
				first_date = excluded.first_date,
				last_date = excluded.last_date,
				data = excluded.data,
				data_no_invest = excluded.data_no_invest`
	istmt, err := txn.PrepareNamed(iquery)
	if err != nil {
		panic(err)
	}
	return istmt
}

func PrepFXTableSt(txn *sqlx.Tx, query string) *sqlx.Stmt {
	stmt, err := txn.Preparex(query)
	if err != nil {
		panic(err)
	}
	return stmt
}
