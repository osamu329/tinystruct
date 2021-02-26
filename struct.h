typedef struct transaction_type {
	char c;
	char s;
	uint16_t n;
} transaction_type_t;

typedef struct query_item {
	char name[10];
} query_item_t;

typedef struct query {
	transaction_type_t transaction_type;
	uint16_t items_n;
	query_item_t item[10];
} query_t;
