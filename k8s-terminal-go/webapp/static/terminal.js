function namespaces() {
    let namespaces = null;

    $.ajax({
        url: "/k8s/namespaces",
        async: false,
        method: "GET",
        success: function (data) {
            if (data.code === "0") {
                namespaces = data.data;
            } else {
                alert(data.errMsg)
            }
        },
        error: function (data) {
            console.log(data);
            alert("Server Interval Error")
        }
    })
    return namespaces;
}

function namespacePods(namespace) {
    let pods = null;

    $.ajax({
        url: "/k8s/namespaces/" + namespace + "/pods",
        async: false,
        method: "GET",
        success: function (data) {
            if (data.code === "0") {
                pods = data.data;
            } else {
                alert(data.errMsg)
            }
        },
        error: function (data) {
            console.log(data);
            alert("Server Interval Error")
        }
    })
    return pods;
}

function terminalExec(namespace, pod, container) {
    if (namespace == null || pod == null || container == null) {
        alert("命名空间、Pod 名称、容器名称不能为空")
    } else {
        if (globalSock !== null) {
            globalSock.close()
        }

        let term = new Terminal({
            fontSize: 14,
            fontFamily: 'Consolas, "Courier New", monospace',
            bellStyle: 'sound',
            cursorBlink: true,
            cols: 150,  // 8
            rows: 50,  // 1.6
            convertEol: false,
            termName: 'xterm',
            cursorStyle: 'block',
            drawBoldTextInBrightColors: true,
            enableBold: true,
            experimentalCharAtlas: 'static',
            fontWeight: 'normal',
            fontWeightBold: 'bold',
            lineHeight: 1.0,
            letterSpacing: 0,
            scrollback: 1000,
            screenKeys: false,
            screenReaderMode: false,
            debug: false,
            macOptionIsMeta: false,
            macOptionClickForcesSelection: false,
            cancelEvents: false,
            disableStdin: false,
            useFlowControl: false,
            allowTransparency: false,
            tabStopWidth: 8,
            theme: undefined,
            rendererType: 'canvas',
            windowsMode: false
        });
        $("#terminal").html("")
        $("#terminal").show()
        term.open(document.getElementById('terminal'));

        let sock = new SockJS(window.location.origin + '/terminal/exec?namespace=' + namespace + "&pod=" +
            pod + "&container=" + container)
        globalSock = sock

        sock.onopen = function () {
            console.log('connection open');
        };
        sock.onmessage = function (e) {
            const msg = JSON.parse(e.data)
            term.write(msg.Data)
        };
        sock.onclose = function () {
            console.log('connection closed');
        };


        term.on('data', function (data) {
            sock.send(
                JSON.stringify({
                    Op: 'stdin',
                    Data: data,
                })
            )
        });
    }
}