# go-redcat
Redcat POS Ploygon API

client := goredcat.NewClient( "https://myurl.com" )

success, err := client.AccessToken( "username", "password")

if success {
  //request a report
  request := goredcat.ReportRequest{
		Distinct: true,
		Start:    0,
		Limit:    100,
	}

	request.Constraints.Operator = "and"

	request.AddField(goredcat.FieldHeaderID)
	request.AddField(goredcat.FieldLineID)
	request.AddField(goredcat.FieldCategoryID)
	request.AddField(goredcat.FieldCategoryName)
	request.AddField(goredcat.FieldPrice)
	request.AddField(goredcat.FieldQty)
	request.AddField(goredcat.FieldGST)
	request.AddField(goredcat.FieldAmount)

	request.AddConstraint(goredcat.FieldHeaderID, ">", 2000)
}