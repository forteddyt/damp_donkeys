import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-registration-complete',
  templateUrl: './registration-complete.component.html',
  styleUrls: ['./registration-complete.component.css']
})
export class RegistrationCompleteComponent implements OnInit {
  registered = "You are now registered for your interview!"
  message = "Thank you for using the CSRS Interviewing Check-In System!"
  goodluck = "Good luck on your interview!"

  constructor(private router: Router) { }

  ngOnInit() {
    setTimeout(() => {
      this.router.navigate(['../']);
  }, 5000);  //5s
  }

}
