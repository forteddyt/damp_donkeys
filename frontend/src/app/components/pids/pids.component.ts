import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-pids',
  templateUrl: './pids.component.html',
  styleUrls: ['./pids.component.css']
})
export class PIDsComponent implements OnInit {
  @Input pid:string;
  constructor() { }

  ngOnInit() {
    this.pid = $scope.pid;
    validateNumber(pid)
  }

  validateNumber(pid:string) {
    var obj = document.getComponentById("pid_input");
    var start = obj.value.indexOf('9');
    var retVal = true;
    if (start != -1 && obj.value.length >= 9) {
        this.pid = obj.value.substring(start, start + 9);
    } else {
        alert('The value you have entered is invalid.');
        retVal = false;
    }
    return retVal;
}

}
