angular.module('ngFilter', []).filter('reqFilter', function() {
    return function(items, protosType, pattern) {
        var result = [];
        if (!protosType && !pattern)
            return items;

        function getMatchFunction() {
            return function(item) {
                if (protosType && item.type != protosType) {
                    return false;
                }

                if( pattern && item.query.toLowerCase().indexOf(pattern.toLowerCase()) < 0) {
                    return false
                }

                return true;
            };
        };
        var matchFunc = getMatchFunction();
        for (var i = 0; i < items.length; i++) {
            var item = items[i];
            if (matchFunc(item)) {
                result.push(item);
            }
        }
        return result;
    };
});
var app = angular.module('netgraph', ['angular-websocket', 'ngFilter'])
app.factory('netdata', function($websocket) {
    var dataStream = $websocket("ws://" + location.host + "/data");
    var streams = {};
    var reqs = [];
    dataStream.onMessage(function(message) {
        var e = JSON.parse(message.data);
        if (!e) {
            return;
        }
        e.timestamp = e['@timestamp'];
        e.request = $.trim(e.request);
        e.response = $.trim(e.response);
        if (e.type == 'http') {
            e.status = e.http.code;
            e.query = e.request_uri;

            var headers = [];
            var headerArr = e.request.split("\n");
            for (var i=0; i < headerArr.length; i++) {
                var arr = headerArr[i].split(':');
                headers.push({'name': arr[0], 'value': arr[1]});
            }
            e.http.request_headers = headers;

            headers = [];
            headerArr = e.response.split("\n");
            for (var i=0; i < headerArr.length; i++) {
                var arr = headerArr[i].split(':');
                headers.push({'name': arr[0], 'value': arr[1]});
            }
            e.http.response_headers = headers;

            e.response = e.raw;

        }
        reqs.push(e);
    });
    var data = {
        reqs: reqs,
        streams: streams,
        sync: function() {
            dataStream.send("sync");
        }
    };
    return data;
})
app.controller('HttpListCtrl', function ($scope, netdata) {
    $scope.reqs = netdata.reqs;
    $scope.showDetail = function($event, req) {
        $scope.selectedReq = req;
        var tr = $event.currentTarget;
        if ($scope.selectedRow) {
            $($scope.selectedRow).attr("style", "");
        }
        $scope.selectedRow = tr;
        $(tr).attr("style", "background-color: lightgreen");
    }
    $scope.getHost = function(req) {
        for (var i = 0; i < req.Headers.length; ++i) {
            var h = req.Headers[i];
            if (h.Name == "Host") {
                return h.Value;
            }
        }
        return null;
    }
    $scope.selectedRow = null;
    $scope.filterType = "Uri";
    $scope.order = "Timestamp";
    netdata.sync();
})
