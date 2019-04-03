import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-registration-complete',
  templateUrl: './registration-complete.component.html',
  styleUrls: ['./registration-complete.component.css']
})
export class RegistrationCompleteComponent implements OnInit {
  registered = "You are now registered for your interview!"
  message = "Thank you for using the CSRS Interviewing check-in system!"
  constructor() { }

  ngOnInit() {
  }

}
