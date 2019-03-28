import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SwipeComponent } from './components/swipe/swipe.component';
import { SelectCompanyComponent } from './components/select-company/select-company.component';
import { ReviewScreenComponent } from './components/review-screen/review-screen.component';
import { RegistrationCompleteComponent } from './components/registration-complete/registration-complete.component';

@NgModule({
  declarations: [
    AppComponent,
    SwipeComponent,
    SelectCompanyComponent,
    ReviewScreenComponent,
    RegistrationCompleteComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
