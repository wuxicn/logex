// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file defines Goid() function, which is used to get the id of current
// goroutine. More details about this function are availeble in the runtime.c
// file of golang source code.

#include <runtime.h>

void ·goid(int32 ret) {
    ret = g->goid;
    FLUSH(&ret);
}
