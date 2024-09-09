
const App = () => {
    const [isValid, setIsValid] = useState(false);
    const [logs, setLogs] = useState([]);
    const [message, setMessage] = useState(INVALID_MESSAGE);
    const [hostname, setHostname] = useState("");

    useEffect(async () => {
        const referrer = new URL(document.referrer);
        // const referrer = new URL("https://example.com")

        setHostname(referrer.hostname)
  
        await setupNotification(referrer.hostname)

        const {v, data} = await fetchAndValidate(referrer.hostname)
        setLogs(data)
        setIsValid(v)
        setMessage(v ? VALID_MESSAGE : INVALID_MESSAGE)

        if (v) {
            setTimeout(() => {
                if (localStorage.stopAutoRedirect == 'true') {
                    setMessage(VALID_AND_STOP_MESSAGE)
                    return
                }

                window.location = document.referrer

            }, 5000)
        }
    }, []);

    return (
        <div>
            <h1>RA-WEBs: TEE Verification Service</h1>
            <h2>Attestation Result</h2>

            <h3>Result: </h3>
            <bold>{isValid ? 'valid' : 'invalid'}</bold>

            <h3>Message: </h3>
            <p>{message}</p>

            <input
                type="checkbox"
                name="subscribe"
                defaultChecked={localStorage.stopAutoRedirect == 'true'}
                onChange={e => {
                    localStorage.stopAutoRedirect = e.target.checked
                }}
            />
            <label for="autoredirect">Do not redirect back automatically</label>


            <br />

            <h2>Logs</h2>
            <TableCompornent logs={logs} />

            <br />
            <button onClick={e => window.location = "https://crt.sh/?q=" + hostname}>See certificate logs (on crt.sh)</button>
        </div>
    );
}