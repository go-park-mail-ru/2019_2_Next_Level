package mailpicker

import (
	"2019_2_Next_Level/internal/MailPicker/log"
	"2019_2_Next_Level/internal/MailPicker/workers"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/tests/mock"
	"2019_2_Next_Level/tests/mock/MailPicker"
	"2019_2_Next_Level/tests/mock/postinterface"
	"context"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
	"time"
)

func init() {
	log.SetLogger(&mock.MockLog{})
}

func Test(t *testing.T) {
	ctx1, finish1 := context.WithCancel(context.Background())
	ctx2, finish2 := context.WithCancel(context.Background())
	ctx3, finish3 := context.WithCancel(context.Background())
	errorChan := make(chan error, 3)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chan1 := make(chan interface{}, 1)
	mockIncoming := postinterface.NewMockIPostInterface(mockCtrl)
	mockRepo := MailPicker.NewMockRepository(mockCtrl)

	picker := workers.NewMailPicker(errorChan, mockIncoming, mockRepo.UserExists)
	cleaner := workers.NewMailCleanup(errorChan)
	chan2 := make(chan model.Email, 1)
	saver := workers.NewMailSaver(errorChan, mockRepo.AddEmail)

	wg := sync.WaitGroup{}
	wg.Add(3)


	returned := post.Email{From:"test", To:"tester", Body:`Received: from mxback18o.mail.yandex.net (mxback18o.mail.yandex.net [IPv6:2a02:6b8:0:1a2d::69])
	by forward102o.mail.yandex.net (Yandex) with ESMTP id 1237E6680F6C
	for <aaa@nlmail.ddns.net>; Sun,  3 Nov 2019 20:54:18 +0300 (MSK)
Received: from localhost (localhost [::1])
	by mxback18o.mail.yandex.net (nwsmtp/Yandex) with ESMTP id rQcGrpkaUN-sHJCr0Vj;
	Sun, 03 Nov 2019 20:54:17 +0300
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed; d=yandex.ru; s=mail; t=1572803657;
	bh=MFjd9jm/+RmHHzeuhW3/etKu87Zh2HmHjtt5+ex1oOI=;
	h=Message-Id:Date:Subject:To:From;
	b=Z2iNN3QnvZzW3K+hne6wT89mxUxKZhU+dtYQBiPUAYYsrDRpwkFI/CT2UMPLNKBOt
	 dTgAgzR+up/ExMuUVR6GSSSj/WEmbWfOyBKySVVlbH1ScDkmwISInwRwjy7prpF8pG
	 C7qOmdGREJXRAnzeucejyEbbGP/Lt+TaqEw+ObEo=
Authentication-Results: mxback18o.mail.yandex.net; dkim=pass header.i=@yandex.ru
Received: by sas8-7ec005b03c91.qloud-c.yandex.net with HTTP;
	Sun, 03 Nov 2019 20:54:17 +0300
From: Andrey K. <andreykochnov@yandex.ru>
To: aaa <aaa@nlmail.ddns.net>
Subject: Test
MIME-Version: 1.0
X-Mailer: Yamail [ http://yandex.ru ] 5.0
Date: Sun, 03 Nov 2019 20:54:17 +0300
Message-Id: <57313241572803657@sas8-7ec005b03c91.qloud-c.yandex.net>
Content-Transfer-Encoding: 7bit
Content-Type: text/html

<div>Hellp</div>`}
	mockIncoming.EXPECT().Get().Return(returned, nil).Times(1)
	mockRepo.EXPECT().UserExists("tester").Return(true).Times(1)
	mockRepo.EXPECT().AddEmail(gomock.Any()).Return(nil).Times(1)

	timer := time.NewTimer(100 * time.Millisecond)
	go picker.Run(&wg, ctx1, chan1)
	select {
	case res := <-chan1:
		finish1()
		chan1<-res
		break
	case <-timer.C:
		t.Errorf("Timeout while waiting for picker result")
		break
	}
	//return
	go cleaner.Run(&wg, ctx2, chan1, chan2)
	timer = time.NewTimer(100 * time.Millisecond)
	select {
	case res := <-chan2:
		finish2()
		chan2<-res
		break
	case <-timer.C:
		t.Errorf("Timeout while waiting for cleaner result")
		break
	}
	go saver.Run(&wg, ctx3, chan2)
	timer = time.NewTimer(10 * time.Millisecond)
	select {
	case <-timer.C:
		finish3()
		break
	}



}