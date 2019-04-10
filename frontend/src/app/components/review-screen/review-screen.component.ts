import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-review-screen',
  templateUrl: './review-screen.component.html',
  styleUrls: ['./review-screen.component.css']
})
export class ReviewScreenComponent implements OnInit {
  stateData

  constructor(private router: Router) {
  	this.stateData = this.router.getCurrentNavigation().extras.state;
    if(this.stateData == null || this.stateData.pid == null || this.stateData.companyName == null){
      console.log("Invalid state has been passed (or not passed at all); should redirect to '/'")
    }else{
    	console.log(this.stateData)
    }
  }

  ngOnInit() {
  	console.log()
  }

}
