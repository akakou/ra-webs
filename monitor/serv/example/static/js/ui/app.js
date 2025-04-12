
const App = () => {
    const [isValid, setIsValid] = useState(false);
    const [logs, setLogs] = useState([]);
    const [message, setMessage] = useState(INVALID_MESSAGE);
    const [hostname, setHostname] = useState("");

    useEffect(async () => {
        const referrer = new URL("http://localhost:8000");
        // const referrer = new URL("https://example.com")

        setHostname(referrer.hostname)
  
        await setupNotification()

        const resp = await axios.get(`/api/ta`)
        console.log(resp)
    
        const logs = resp.data
        const v = checkValidity(logs)

        var message = v ? VALID_MESSAGE : INVALID_MESSAGE

        setLogs(logs)
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