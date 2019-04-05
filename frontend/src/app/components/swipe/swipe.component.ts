import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router, ActivatedRoute, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-swipe',
  templateUrl: './swipe.component.html',
  styleUrls: ['./swipe.component.css']
})
export class SwipeComponent implements OnInit {
  title = "Enter PID"
  PIDInput = "Enter your PID"
  pid = ''
  student_info
  constructor(private http: HttpClient, private router: Router, private route: ActivatedRoute) { }

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
    event.target.value = this.pid;
    if (this.pid.length == 9 && /^\d+$/.test(this.pid)) {
      var student = this.http.get("http://csrcint.cs.vt.edu:8080/get_student?VT_ID=" + this.pid).subscribe((res) => {
        console.log(res);
        this.student_info = res;
      });
      let navigationExtras: NavigationExtras = {
        queryParams: {
          student: this.student_info.dispName,
          class: this.student_info.class,
          major: this.student_info.major
        },
        skipLocationChange: true
      }
      this.router.navigate(['../select'], navigationExtras);
    }
  }
}
