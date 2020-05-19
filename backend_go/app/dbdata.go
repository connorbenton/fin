package app

import (
	"time"
)

type category struct {
	id   int				`json:"id"`
	topCategory	string		`json:"top_category"`		
	subCategory string		`json:"sub_category"`
	excludeFromAnalysis int	`json:"exclude_from_analysis"`
	createdAt time.Time 	`json:"created_at"`
	updatedAt time.Time 	`json:"updated_at"`
}

type account struct {
	Id   		int			`json:"id"`
	Name		string		`json:"name"`		
	Institution string		`json:"institution"`
	Account_id 	string		`json:"account_id"`
	Item_id 	string		`json:"item_id"`
	Type 		string 		`json:"type"`
	Subtype 	string		`json:"subtype"`
	Balance 	float32		`json:"balance"`
	Limit 		float32 	`json:"limit"`
	Available 	float32 	`json:"available"`
	Currency 	string 		`json:"currency"`
	Provider	string 		`json:"provider"`
	CreatedAt time.Time 	`json:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at"`
}

type itemToken struct {
	id   			int			`json:"id"`
	institution		string		`json:"institution"`		
	access_token	string		`json:"access_token"`
	item_id		 	string		`json:"item_id"`
	provider		string		`json:"provider"`
	interactive		bool		`json:"interactive"`
	needsReLogin	bool		`json:"needs_re_login"`
	lastRefresh 				time.Time 	`json:"last_refresh"`
	nextRefreshPossible 		time.Time 	`json:"last_refresh"`
	lastDownloadedTransactions 	time.Time 	`json:"last_refresh"`
	createdAt time.Time 	`json:"created_at"`
	updatedAt time.Time 	`json:"updated_at"`
}