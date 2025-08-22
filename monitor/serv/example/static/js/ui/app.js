
const App = () => {
    const [isValid, setIsValid] = useState(false);
    const [logs, setLogs] = useState([]);
    const [message, setMessage] = useState(INVALID_MESSAGE);
    const [hostname, setHostname] = useState("");

    useEffect(async () => {
        const resp1 = await axios.get(`/api/config`)
        console.log(resp1)

        setHostname(resp1.data.taDomain)
  
        const resp2 = await axios.get(`/api/ta`)
        console.log(resp2)
    
        const logs = resp2.data
        const v = checkValidities(logs)

        var message = v ? VALID_MESSAGE : INVALID_MESSAGE

        setLogs(logs)
        setIsValid(v)
        setMessage(message)

    }, []);

    return (
        <div>
            <h1>RA-WEBs: TEE Monitoring Service</h1>
            <h2>Attestation Result</h2>

            <button onClick={async e => {
                await setupNotification()
            }}>Set up Notification</button>

            <h3>TA Domain: </h3>
            <a href={"https://" + hostname}>{hostname}</a>

            <h3>Result: </h3>
            <bold>{isValid ? 'valid' : 'invalid'}</bold>

            <h3>Message: </h3>
            <p>{message}</p>

            <h2>Logs</h2>
            <TableCompornent logs={logs} />

            <br />
            <button onClick={e => window.location = "https://crt.sh/?q=" + hostname}>See certificate logs (on crt.sh)</button>
        </div>
    );
}