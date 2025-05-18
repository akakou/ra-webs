const TableCompornent = ({ logs }) => {
    const gitRepo = s => {
        if (s.edges.evidence_log) {
            let href = `${s.edges.evidence_log.repository}/tree/${s.edges.evidence_log.commit_id}`
            return <button 
                    onClick={
                        e => window.location = href
                    }
                    >Go GitHub Repo
                </button>
        } else {
            return <p>No Repo</p>
        }
    } 
    
    const crtshLinks = server => server.edges.ct_log.map(
        (ctl, index) =>
            <div>
                <a href={"https://crt.sh?id="+ctl.monitor_log_id}>
                    {ctl.monitor_log_id}  
                </a>
                <br/>
            </div>
        )
    const uniqueId = (server) => server.edges.evidence_log ? server.edges.evidence_log.unique_id : "-"
    const violated = (server) =>  
        (!checkValidity(server)).toString()
    

    console.log(logs)
    const rows = logs.map((server, index) =>
        <tr key={index}>
            <td>{server.id}</td>
            <td>{uniqueId(server)}</td>
            <td>{crtshLinks(server)}</td>
            <td>{gitRepo(server)}</td>
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
