//Define an angular module for our app
var checkIn = angular.module('check', []);

//Define Routing for app
//Uri /AddNewOrder -> template add_order.html and Controller AddOrderController
//Uri /ShowOrders -> template show_orders.html and Controller AddOrderController
checkIn.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/swipe', {
		templateUrl: './swipe.html',
		controller: 'SwipeController'
	}).
      when('/ShowOrders', {
		templateUrl: 'templates/show_orders.html',
		controller: 'ShowOrdersController'
      }).
      otherwise({
		redirectTo: '/swipe'
      });
}]);


sampleApp.controller('SwipeController', function($scope) {
	
	$scope.message = 'This is Add new order screen';
	
});


sampleApp.controller('ShowOrdersController', function($scope) {

	$scope.message = 'This is Show orders screen';

});
