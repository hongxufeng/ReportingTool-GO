package model

type ColumnConfig struct {
	Tag             string
	Text            string
	HasBtn          bool
	HasDateformat   bool
	HasDefaultValue bool
	HasFormatter    bool
	HasFormatterR   bool
	//HasLinkTo bool
	HasNavName   bool
	HasPrecision bool
	//HasRegex bool
	HasSearchType    bool
	HasSelectorMulti bool
	HasTimeTransfer  bool
	HasSearchBtnIcon bool
	HasSearchBtnFunc bool
	HasSelectorFunc  bool
	HasSelectorText  bool
	//IsInPercentageform bool
	IsInSelector  bool
	Search4Admin  bool
	ISCheckBox    bool
	IsSearchAdv   bool
	HasPower      bool
	Power         uint8
	SearchBtnIcon string
	SearchBtnFunc string
	BtnIcon       string
	BtnFunc       string
	ColumnName    string
	DateFormat    string
	DefaultValue  string
	Formatter     string
	FormatterR    string
	//LinkTo string
	NavName      string
	TimeTransfer string
	Precision    string
	//RegexPattern string
	//RegexReplacement string
	SearchType       string
	Selector         string
	SelectorFunc     string
	SelectorFuncAgrs string
	SelectorText     string
	Visibility       string
	//Passedcol []string
}
type TableConfig struct {
	Name            string
	DefaultOrder    string
	HasDefaultOrder bool
	Excel           string //生成excel
	HasExcel        bool
	//AdminName string
	//HasAdminName bool
	Power    uint8
	HasPower bool
}
