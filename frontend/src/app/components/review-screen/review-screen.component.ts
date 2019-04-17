import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
@Component({
  selector: 'app-review-screen',
  templateUrl: './review-screen.component.html',
  styleUrls: ['./review-screen.component.css']
})
export class ReviewScreenComponent implements OnInit {
  stateData

  constructor(private router: Router, private http: HttpClient) {
  	this.stateData = this.router.getCurrentNavigation().extras.state;
    if(this.stateData == null || this.stateData.pid == null || this.stateData.companyName == null){
      //console.log("Invalid state has been passed (or not passed at all); should redirect to '/'");
      this.router.navigate(['/']);
    }
  }

  ngOnInit() { }

  submitInterview(event: any) {
    //console.log("Submit the Interview Details");
    this.http.put("https://csrcint.cs.vt.edu/api/interview_check_in?",
      {
        "company_name": this.stateData.companyName,
        "display_name": this.stateData.student,
        "major": this.stateData.major,
        "class": this.stateData.class,
        "VT_ID": this.stateData.pid
      });
    this.router.navigate(['/complete']);
  }
}
