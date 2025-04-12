const TableCompornent = ({ logs }) => {
    const gitRepo = s => `${s.edges.at_log[0].repository}/tree/${s.edges.at_log[0].commit_id}`
    console.log(logs)
    const rows = logs.map((server, index) =>
        <tr key={index}>
            <td>{server.id}</td>
            <td>{server.edges.at_log[0].unique_id}</td>
            <td>
                <a href={"https://crt.sh?id="+server.edges.ct_log[0].monitor_log_id}>
                {server.edges.ct_log[0].monitor_log_id}
                </a>
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
