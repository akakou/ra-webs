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
    
    const publicKey = server => server.public_key ? server.public_key : "-"
    const uniqueId = server => server.edges.evidence_log ? server.edges.evidence_log.unique_id : "-"
    const evidence = server => server.edges.evidence_log ? server.edges.evidence_log.evidence : "-"
    const violated = server =>  
        (!checkValidity(server)).toString()
    

    console.log(logs)
    const rows = logs.map((server, index) =>
        <tr key={index}>
            <td>{server.id}</td>
            <td>{uniqueId(server)}</td>
            <td>
                <textarea readOnly defaultValue={publicKey(server)} style={{width: "200px", height: "200px", resize: "none"}}></textarea>
            </td>
            <td>
                <textarea readOnly defaultValue={evidence(server)} style={{width: "200px", height: "200px", resize: "none"}}></textarea>
            </td>
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
                    <th>Public Key</th>
                    <th>Evidence</th>
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
