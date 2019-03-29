import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-swipe',
  templateUrl: './swipe.component.html',
  styleUrls: ['./swipe.component.css']
})
export class SwipeComponent implements OnInit {
  title = "Enter PID"
  PIDInput = "Enter your PID"
  pid = ''
  constructor() { }

  ngOnInit() {
    document.getElementById("pid_input").focus();
  }

  onKey(event: any) { // without type info
    var input = event.target.value;
    if (input.length < 9) {
      this.pid = event.target.value;
    }
    else {
      var start = input.indexOf('9');
      if (start != -1) {
        this.pid = input.substring(start, start + 9);
      } 
    }
    if (this.pid.length == 9 && /^\d+$/.test(this.pid)) {
      //API Call here and move onto next page
    }
    event.target.value = this.pid;
  }
}
