package main

import (
	"./lib/utils"
	"./lib/bayesian"
	"fmt"
	"io/ioutil"
	// "os"
	// "reflect"
	"strings"
)

type WordDict map[string]bool
type FeatMat [][]int
type LabelMat []int

// BuildDictionary creates dictionary from all the emails in directory
func BuildDictionary(dir string) ([]utils.WordDict, error) {
	var err error
	// Read the file names and sorts them.
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}
	
	// Slice to hold all the words in the emails
	goodwordlist := []string{}
	spamwordlist := []string{}

	// Collecting all words from those emails
	for _, email := range emailList {
		// This labels the data for training purposes for spam dataset
		if strings.Contains(email.Name(), "spms") {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return nil, fmt.Errorf("File opening failed %s", err)
		}
		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")
		for i, line := range dat{
			// Body of email is only 3rd line of text file
			if i == 2{
				words := strings.Split(line, " ")
				spamwordlist = append(spamwordlist, words...)
			}
		}
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return nil, fmt.Errorf("File opening failed %s", err)
		}
		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")
		for i, line := range dat{
			// Body of email is only 3rd line of text file
			if i == 2{
				words := strings.Split(line, " ")
				goodwordlist = append(goodwordlist, words...)
			}
		}
	}
	}

	// STATS: Wordcount -> 138777
	fmt.Println(len(goodwordlist))
	fmt.Println(len(spamwordlist))
	
	// We now have the dictionary of words, which may have duplicate entries
	goodworddict := utils.Set(goodwordlist) // Duplicates removed.
	spamworddict := utils.Set(spamwordlist) // Duplicates removed.
	
	// STATS: Wordcount -> 13397
	fmt.Println(len(goodworddict))
	fmt.Println(len(spamworddict))

	// Removes punctuations and non-alphabets
	for word := range goodworddict{
		if len(word) == 1 || !utils.IsAlpha(word) {
			delete(goodworddict, word)
		}
	}

	for word := range spamworddict{
		if len(word) == 1 || !utils.IsAlpha(word) {
			delete(spamworddict, word)
		}
	}

	// STATS: Wordcount -> 11793
	fmt.Println(len(goodworddict))
	fmt.Println(len(spamworddict))

	return []utils.WordDict{goodworddict, spamworddict}, nil
}

// BuildFeatures returns the feature matrix
func BuildFeatures(dir string, dictionary utils.WordDict) (FeatMat, error) {
	// Read the file names
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}

	// Matrix to have features
	// Renders a matrix like this featMatrix => [len(emailList)][len(dictionary)]int
	featMatrix := make([][]int, len(emailList))
	for i := range featMatrix {
		featMatrix[i] = make([]int, len(dictionary))
	}

	// Collecting the number of occurences of each of the words in the emails.
	for emailI, email := range emailList {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, email.Name()))
		if err != nil {
			return nil, fmt.Errorf("File opening failed %s", err)
		}
		// Breaks the email into lines.
		dat := strings.Split(string(data), "\n")
		for lineI, line := range dat{
			// Body of email is only 3rd line of text file
			if lineI == 2{
				words := strings.Split(line, " ")
				for word, wordI := range dictionary{
					featMatrix[emailI][wordI] = utils.Count(words, word)
				}
			}
		}
	}
	return featMatrix, nil
}

func BuildLabels(dir string) (LabelMat, error){
	// Read the file names
	emailList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory not present %s", err)
	}

	// Label vector
	labelMat := make([]int, len(emailList))

	for i, email := range emailList{
		if strings.Contains(email.Name(), "spms") {
			labelMat[i] = 1 
		} else {
			labelMat[i] = 0
		}
	}

	return labelMat, nil
}

func main() {
	trainDir := "dataset/train_data"
	
	fmt.Println("1. Building dictionary")
	dict, err := BuildDictionary(trainDir)
	utils.Check(err)

	// fmt.Println("2. Building training features and labels")
	// featTrain, err := BuildFeatures(trainDir, dict)
	// utils.Check(err)
	// labelTrain, err := BuildLabels(trainDir)
	// utils.Check(err)

	// fmt.Println("3. Training the Classifier")
	// featTrain, err := BuildFeatures(trainDir, dict)
	// utils.Check(err)

	var (
		Fam bayesian.Class = "Fam"   // The good ones
		Spam bayesian.Class = "Spam" // The bad ones
	)
	classifier := bayesian.NewClassifierTfIdf(Spam, Fam)

	famMails := utils.MapToArr(dict[0])
	spamMails := utils.MapToArr(dict[1])

	classifier.Learn(famMails, Fam)
	classifier.Learn(spamMails, Spam)

	classifier.ConvertTermsFreqToTfIdf()

	// scores, likely, _ := classifier.LogScores(
	// 	[]string{}
	// )
}