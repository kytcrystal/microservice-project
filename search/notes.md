# Steps to get rabbit mq running

1. Run: docker compose start 
2. Check that the management console is up and running
   - Go to http://localhost:15672 the default user name is guest and password is also guest
3. (if not already done) Run go get github.com/rabbitmq/amqp091-go
    - this will add rabbit mq library as dependency of our program in go.mod (and install required stuff)
4. 

go get github.com/rabbitmq/amqp091-go


## Resources:

- https://www.rabbitmq.com/tutorials/tutorial-one-go.html
