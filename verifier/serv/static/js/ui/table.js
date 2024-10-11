const TableCompornent = ({ logs }) => {
    const gitRepo = s => `${s.edges.code.repository}/tree/${s.edges.code.commit_id}`

    const rows = logs.map((server, index) =>
        <tr key={index}>
            <td>{server.id}</td>
            <td>{server.domain}</td>
            <td>{server.edges.code.unique_id}</td>
            <td>
                <a href={"https://crt.sh?id="+server.monitor_log_id}>
                {server.monitor_log_id}
                </a>
            </td>
            <td>
                <a href={gitRepo(server)}>{gitRepo(server)}</a>
            </td>
            <td>{(!!server.is_active).toString()}</td>
            <td>{(server.edges.violation.length == 0).toString()}</td>
        </tr>
    );

    return (
        <table>
            <thead>
                <tr>
                    <th>Index</th>
                    <th>Domain</th>
                    <th>Unique ID</th>
                    <th>crt.sh ID</th>
                    <th>Git Repository</th>
                    <th>Activated</th>
                    <th>Violated</th>
                </tr>
            </thead>
            <tbody>
                {rows}
            </tbody>
        </table>
    )
}
