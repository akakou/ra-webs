self.addEventListener('push', function (event) {
    console.log(event)
    const text = event.data.text();
    const lines = text.split('\n')
    const title = lines[0]
    const message = lines.slice(1).join('\n')

    event.waitUntil(
        self.registration.showNotification(title, {
            body: message,
            tag: 'push-notification-tag'
        })
    );
});
