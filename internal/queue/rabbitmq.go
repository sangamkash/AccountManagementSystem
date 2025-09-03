package queue

import (
    "context"
    "encoding/json"
    "log"
    "sync"

    "github.com/streadway/amqp"
)

// TransactionMessage defines the payload
type TransactionMessage struct {
    IdempotencyKey string  `json:"idempotency_key"`
    AccountID      int64   `json:"account_id"`
    Type           string  `json:"type"`
    Amount         float64 `json:"amount"`
}

type Rabbit struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   string
    mu      sync.Mutex
}

func NewRabbit(url, queueName string) (*Rabbit, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }
    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }
    _, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, err
    }
    return &Rabbit{conn: conn, channel: ch, queue: queueName}, nil
}

func (r *Rabbit) Publish(ctx context.Context, msg TransactionMessage) error {
    b, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    // make publish thread-safe on the channel
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.channel.Publish("", r.queue, false, false, amqp.Publishing{
        ContentType: "application/json",
        Body:        b,
        DeliveryMode: amqp.Persistent,
    })
}

func (r *Rabbit) Consume(ctx context.Context, handler func(context.Context, TransactionMessage) error) error {
    msgs, err := r.channel.Consume(r.queue, "", false, false, false, false, nil)
    if err != nil {
        return err
    }

    for d := range msgs {
        var m TransactionMessage
        if err := json.Unmarshal(d.Body, &m); err != nil {
            log.Printf("invalid message: %v", err)
            d.Nack(false, false)
            continue
        }
        if err := handler(ctx, m); err != nil {
            log.Printf("handler error: %v", err)
            d.Nack(false, true) // requeue
            continue
        }
        d.Ack(false)
    }
    return nil
}

func (r *Rabbit) Close() {
    if r.channel != nil {
        r.channel.Close()
    }
    if r.conn != nil {
        r.conn.Close()
    }
}
