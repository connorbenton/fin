package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name" db:"name"`
	Institution        string          `json:"institution" db:"institution"`
	IgnoreTransactions string          `json:"ignore_transactions" db:"ignore_transactions"`
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
	ID                         int       `json:"id"`
	Institution                string    `json:"institution" db:"institution"`
	AccessToken                string    `json:"access_token" db:"access_token"`
	ItemID                     string    `json:"item_id" db:"item_id"`
	Provider                   string    `json:"provider" db:"provider"`
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
	ExcludeFromAnalysis int       `json:"exclude_from_analysis" db:"exclude_from_analysis"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID                  int       `json:"id"`
	Date                time.Time `json:"date" db:"date"`
	TransactionID       string    `json:"top_category" db:"top_category"`
	Description         string    `json:"sub_category" db:"sub_category"`
	OriginalDescription int       `json:"exclude_from_analysis" db:"exclude_from_analysis"`
	Amount              time.Time `json:"created_at" db:"created_at"`
	NormalizedAmount    time.Time `json:"updated_at" db:"updated_at"`
	TransactionType     time.Time `json:"updated_at" db:"updated_at"`
}

type EcbFX struct {
	Currencies []struct {
		SeriesKey struct {
			Text   string `xml:",chardata"`
			Values []struct {
				ID    string `xml:"id,attr"`
				Value string `xml:"value,attr"`
			} `xml:"Value"`
		} `xml:"SeriesKey"`
		Rates []struct {
			Date struct {
				Value string `xml:"value,attr"`
			} `xml:"ObsDimension"`
			Rate struct {
				Value string `xml:"value,attr"`
			} `xml:"ObsValue"`
		} `xml:"Obs"`
	} `xml:"DataSet>Series"`
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
		CategorizationConfidence int             `json:"categorization_confidence"`
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
