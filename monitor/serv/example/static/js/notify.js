const setupNotification = async () => {
    removeAllNotifcation()
    const keyResp = await axios.get(`/api/config/subscription`)
    console.log(keyResp)

    const permission = await Notification.requestPermission()
    if (permission !== 'granted') {
        console.error("Notification is not granted")
        return
    }

    const registrations = await navigator.serviceWorker.getRegistrations()
    if (registrations.length > 0) {
        console.log("Service worker already registered")
        return 
    }

    let subscription
    try {
        await navigator.serviceWorker.register('/static/js/sw/service-worker.js', {scope: "/"})
        const sw = await navigator.serviceWorker.ready;

        subscription = await sw.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: keyResp.data.vapid_public_key,
        });
    } catch (e) {
        console.error("Service worker registration failed: ", e)
        return
    }

    const subscribeResp = await axios.post(`/api/subscription`, {
        subscription: subscription.toJSON()
    })

    console.log("registered: ", subscribeResp)

}

// this is for debug
const removeAllNotifcation = async () => {
    const registrations = await navigator.serviceWorker.getRegistrations()
    for (const registration of registrations) {
        await registration.unregister()
    }
    console.log("unregistered")
}