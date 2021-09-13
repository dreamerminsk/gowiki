package tasks


func updateCategories() {
        cats, err := getCategories()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	for _, cat := range cats {
		fmt.Println("Title: ", cat.Title)
		g.DB.Create(cat)
	}
}
