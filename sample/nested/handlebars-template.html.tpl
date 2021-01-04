<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>A simple, clean, and responsive HTML invoice template</title>
</head>

<body>
    <div>
        {{#each list}}
        <tr>
            <td>
                {{listProperty}}
                {{#each nestedList}}
                    {{.}}
                {{/each}}
            </td>
        </tr>
        {{/each}}
    </div>
</body>
</html>