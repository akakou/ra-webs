const TableCompornent = ({ logs }) => {
    const gitRepo = s => `${s.edges.at_log.repository}/tree/${s.edges.at_log.commit_id}`
    const crtshLinks = server => server.edges.ct_log.map(
        (ctl, index) =>
            <div>
                <a href={"https://crt.sh?id="+ctl.monitor_log_id}>
                    {ctl.monitor_log_id}  
                </a>
                <br/>
            </div>
        )
    const uniqueId = (server) => {
        if (!!server.edges.at_log)
            return server.edges.at_log.unique_id
        else 
            return ""
    }
    const violated = (server) =>  
        (checkValidity(server)).toString()
    

    console.log(logs)
    const rows = logs.map((server, index) =>
        <tr key={index}>
            <td>{server.id}</td>
            <td>{uniqueId(server)}</td>
            <td>
            {crtshLinks(server)}
            </td>

            <td>
                <a href={gitRepo(server)}>{gitRepo(server)}</a>
            </td>
            <td>{violated(server)}</td> 
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
