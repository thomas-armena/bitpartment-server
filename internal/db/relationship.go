package db

//Relationship represents how one tenant feels about another tenant
type Relationship struct {
	tableName  struct{} `sql:"relationships"`
	Tenant1ID  int      `sql:"tenant_1_id"`
	Tenant2ID  int      `sql:"tenant_2_id"`
	Romance    int      `sql:"romance,type:smallint"`
	Friendship int      `sql:"romance,type:smallint"`
}
