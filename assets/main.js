const { RSocketClient } = require('rsocket-core');
const { BufferEncoders } = require('rsocket-core');
const { ReactiveSocket, Single, Flowable } = require('rsocket-flowable');

// Создаем клиентское соединение
const client = new RSocketClient({
    transport: require('rsocket-tcp-client').connect({
        host: 'localhost',
        port: 11000,
    }),
    setup: {
        dataMimeType: 'text/plain',
        metadataMimeType: 'text/plain',
    },
    responder: {
        requestResponse: (payload) => {
            console.log(`Received response: ${payload.data}`);
        },
        requestStream: (payload) => {
            const response = new Flowable((subscriber) => {
                // Отправляем элементы потока на сервер и обрабатываем ответы
                // Используйте subscriber.onNext() для отправки элементов
                // И subscriber.onComplete() для завершения потока
            });
            return response;
        },
    },
});

// Устанавливаем соединение с сервером
client.connect().subscribe({
    onComplete: (socket) => {
        // Отправляем запросы на сервер
        socket.requestResponse({
            data: 'Ваш запрос',
            metadata: '',
        }).subscribe({
            onComplete: (response) => {
                console.log(`Response: ${response.data}`);
            },
            onError: (error) => {
                console.error(`Error: ${error}`);
            },
        });

        socket.requestStream({
            data: 'Ваш запрос',
            metadata: '',
        }).subscribe({
            onComplete: (response) => {
                console.log(`Stream Response: ${response.data}`);
            },
            onError: (error) => {
                console.error(`Stream Error: ${error}`);
            },
        });
    },
    onError: (error) => {
        console.error(`Connection Error: ${error}`);
    },
});

// Выполните этот код из фронтенд-приложения
