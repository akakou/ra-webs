
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

        const resp = await axios.get(`/api/ta/${referrer.hostname}`)
        console.log(resp)
    
        const ta = resp.data.ta
        const v = resp.data.is_valid

        var message = v ? VALID_MESSAGE : INVALID_MESSAGE
        message += " " + resp.data.message

        setLogs(ta)
        setIsValid(v)
        setMessage(message)

    }, []);

    return (
        <div>
            <h1>RA-WEBs: TEE Verification Service</h1>
            <h2>Attestation Result</h2>

            <h3>Result: </h3>
            <bold>{isValid ? 'valid' : 'invalid'}</bold>

            <h3>Message: </h3>
            <p>{message}</p>

            <br />

            <h2>Logs</h2>
            <TableCompornent logs={logs} />

            <br />
            <button onClick={e => window.location = "https://crt.sh/?q=" + hostname}>See certificate logs (on crt.sh)</button>
        </div>
    );
}