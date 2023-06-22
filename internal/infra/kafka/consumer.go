package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Consumer struct {
	ConfigMap *ckafka.ConfigMap
	Topics    []string
}
 // Função consumir pro kafka
func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}
 // Método canal que recebe as msg ckafka.Message
func (c *Consumer) Consume(msgChan chan *ckafka.Message) error {
	consumer, err := ckafka.NewConsumer(c.ConfigMap) // Variaveis de concumo
	if err != nil { // Verificação de erros
		panic(err)
	}
	err = consumer.SubscribeTopics(c.Topics, nil) 
	if err != nil { // Erro ao se inscrever
		panic(err)
	}
	for { // Loop infinito lendo msg do topico e enviando pro canal
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			msgChan <- msg
		}
	}
}