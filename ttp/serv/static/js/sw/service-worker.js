self.addEventListener('push', function (event) {
    console.log('Received a push message', event);
    const data = event.data?.json() ?? {};

    const title = data.title || "Something Has Happened";
    const message =
      data.message || "Here's something you might want to check out.";

    event.waitUntil(
        self.registration.showNotification(title, {
            body: message,
            tag: 'push-notification-tag'
        })
    );
});