import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SwipeComponent } from './components/swipe/swipe.component';
import { SelectCompanyComponent } from './components/select-company/select-company.component';
import { ReviewScreenComponent } from './components/review-screen/review-screen.component';
import { RegistrationCompleteComponent } from './components/registration-complete/registration-complete.component';
import { CompanyUniqueComponent } from './components/company-unique/company-unique.component';
import { CompanyWelcomeComponent } from './components/company-welcome/company-welcome.component';
import { AdminComponent } from './components/admin/admin.component';

@NgModule({
  declarations: [
    AppComponent,
    SwipeComponent,
    SelectCompanyComponent,
    ReviewScreenComponent,
    RegistrationCompleteComponent,
    CompanyUniqueComponent,
    CompanyWelcomeComponent,
    AdminComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
