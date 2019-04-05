import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-review-screen',
  templateUrl: './review-screen.component.html',
  styleUrls: ['./review-screen.component.css']
})
export class ReviewScreenComponent implements OnInit {
  name = ""
  title = "Review"
  constructor() { }

  ngOnInit() {
  }

}
