self.addEventListener('push', function (event) {
    console.log('Received a push message', event);
    var title = "";
    var body = "プッシュ通知はこのようにして送られるのです";

    event.waitUntil(
        self.registration.showNotification(title, {
            body: body,
            tag: 'push-notification-tag'
        })
    );
});