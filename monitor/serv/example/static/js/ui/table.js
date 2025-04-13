const TableCompornent = ({ logs }) => {
    const gitRepo = s => `${s.edges.at_log.repository}/tree/${s.edges.at_log.commit_id}`
    const crtshLinks = ctLogs => ctLogs.map(
        (ctLog, index) =>
            <div>
                <a href={"https://crt.sh?id="+ctLog.monitor_log_id}>
                    {ctLog.monitor_log_id}  
                </a>
                <br/>
            </div>
        )
     
    console.log(logs)
    const rows = logs.map((server, index) =>
        <tr key={index}>
            <td>{server.id}</td>
            <td>{server.edges.at_log.unique_id}</td>
            <td>
            {crtshLinks(server.edges.ct_log)}
            </td>

            <td>
                <a href={gitRepo(server)}>{gitRepo(server)}</a>
            </td>
            <td>{(!!server.edges.violation).toString()}</td> 
        </tr>
    );


    return (
        <table>
            <thead>
                <tr>
                    <th>Index</th>
                    <th>Unique ID</th>
                    <th>crt.sh ID</th>
                    <th>Git Repository</th>
                    <th>Violated</th>
                </tr>
            </thead>
            <tbody>
                 {rows} 
            </tbody>
        </table>
    )
}
