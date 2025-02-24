package answerinput

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
)

var (
	// GetInstance retrieves instance of the answer input manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*AnswerInputManager](newAnswerInputManager)
)

// Describes all the available operator definitions.
const (
	minusOperator    = "-"
	plusOperator     = "+"
	multiplyOperator = "*"
)

var (
	// Represents all the available operators.
	operators = []string{plusOperator, minusOperator, multiplyOperator}
)

const (
	// Represents max amount of operators.
	maxOperators = 3

	// Represents max value for each section.
	maxValue = 20
)

// AnswerInputManager represents answer input manager, which prepares questions.
type AnswerInputManager struct {
	generatedQuestion *dto.GeneratedQuestionUnit
}

// UpdateQuestion updates question for answer input.
func (sm *AnswerInputManager) UpdateQuestion() {
	rand.Seed(time.Now().UnixNano())

	number := rand.Intn(maxValue) + 1

	sequence := []string{fmt.Sprintf("%d", number)}

	var operator string

	for i := 0; i < maxOperators-1; i++ {
		operator = operators[rand.Intn(len(operators))]

		number = rand.Intn(maxValue) + 1

		sequence = append(sequence, operator, fmt.Sprintf("%d", number))
	}

	question := strings.Join(sequence, " ")

	expression, err := govaluate.NewEvaluableExpression(question)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	sm.generatedQuestion = &dto.GeneratedQuestionUnit{
		Question: question,
		Answer:   int(result.(float64)),
	}
}

// GetGeneratedQuestion retrieves generated question for answer input.
func (sm *AnswerInputManager) GetGeneratedQuestion() *dto.GeneratedQuestionUnit {
	return sm.generatedQuestion
}

// newAnswerInputManager initializes AnswerInputManager.
func newAnswerInputManager() *AnswerInputManager {
	return new(AnswerInputManager)
}
