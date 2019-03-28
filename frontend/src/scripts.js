var module = angular.module("frontpage", []);

module.controller("formCtrl", ['$scope', function ($scope) {
    $scope.pid = {{pid}}
}]);

function validateNumber() {
    var obj = $scope.pid;
    var start = obj.value.indexOf('9');
    var retVal = true;
    if (start != -1 && obj.value.length >= 9) {
        obj.value = obj.value.substring(start, start + 9);
    } else {
        alert('The value you have entered is invalid.');
        retVal = false;
    }
    return retVal;
}

