package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"text/template"
)

type recipeRequest struct {
	water         float64
	salt          float64
	yeast         float64
	gramsPerPizza float64
	numPizzas     float64
}

type recipe struct {
	Flour         float64 `json:"flour"`
	Water         float64 `json:"water"`
	Salt          float64 `json:"salt"`
	Yeast         float64 `json:"yeast"`
	GramsPerPizza float64 `json:"gramsPerPizza"`
	NumPizzas     float64 `json:"numPizzas"`
}

var defaultRecipe = recipeRequest{
	water:         64.0,
	salt:          2.0,
	yeast:         0.5,
	gramsPerPizza: 250,
	numPizzas:     4}

func main() {

	waterPtr := flag.Float64("water", defaultRecipe.water, "Wanted water percentage")
	saltPtr := flag.Float64("salt", defaultRecipe.salt, "Wanted salt percentage")
	yeastPtr := flag.Float64("yeast", defaultRecipe.yeast, "Wanted water percentage")
	gramsPerPizza := flag.Float64("grams-per-pizza", defaultRecipe.gramsPerPizza, "grams per pizza")
	numPizzas := flag.Float64("num-pizzas", defaultRecipe.numPizzas, "Number of pizzas")

	flag.Parse()

	req := recipeRequest{water: *waterPtr,
		salt:          *saltPtr,
		yeast:         *yeastPtr,
		gramsPerPizza: *gramsPerPizza,
		numPizzas:     *numPizzas}

	recipe := generateRecipe(req)

	printRecipie(recipe)

}

func printRecipie(r recipe) {
	tmpl, _ := template.New("recipe").Funcs(template.FuncMap{
		"d0": func(f float64) string { return fmt.Sprintf("%.0f", f) },
	}).Parse(
		"Flour: {{d0 .Flour}} g\n" +
			"Water: {{d0 .Water}} g\n" +
			"Salt:  {{d0 .Salt}} g \n" +
			"Yeast: {{d0 .Yeast}} g\n\n" +
			"Mix water and yeast\n" +
			"Add salt and almost all flour.\n" +
			"Add rest of flour (if needed).\n" +
			"Knead, divide into {{d0 .NumPizzas}} and refridigrate for 24h\n")

	tmpl.Execute(os.Stdout, r)
}

func toDecimal(percent float64) float64 {
	return percent / 100.0
}

func generateRecipe(req recipeRequest) recipe {
	totalDoughGrams := req.numPizzas * req.gramsPerPizza
	flourGrams := totalDoughGrams / (1 + toDecimal(req.water) + toDecimal(req.salt) + toDecimal(req.yeast))

	return recipe{
		GramsPerPizza: req.gramsPerPizza,
		NumPizzas:     req.numPizzas,
		Salt:          flourGrams * toDecimal(req.salt),
		Water:         flourGrams * toDecimal(req.water),
		Yeast:         flourGrams * toDecimal(req.yeast),
		Flour:         flourGrams}
}

func createRecipeRequest(data url.Values) recipeRequest {
	water := getFloat(data.Get("water"), defaultRecipe.water)
	salt := getFloat(data.Get("salt"), defaultRecipe.salt)
	yeast := getFloat(data.Get("yeast"), defaultRecipe.yeast)
	gramsPerPizza := getFloat(data.Get("gramsPerPizza"), defaultRecipe.gramsPerPizza)
	numPizzas := getFloat(data.Get("numPizzas"), defaultRecipe.numPizzas)
	return recipeRequest{gramsPerPizza: gramsPerPizza, numPizzas: numPizzas, salt: salt, yeast: yeast, water: water}
}

func getFloat(value string, def float64) float64 {

	if s, err := strconv.ParseFloat(value, 64); err == nil {
		return s
	}

	return def
}
