package utils
import (
"net/url"
)





func GetParam(ref, name string)(*string, error) {
    u, err := url.Parse(ref)
    if err != nil {
        return nil, err
    }
    q, err := url.ParseQuery(u.RawQuery)
    if err != nil {
        return nil, err
    }
    return &q[name][0], nil
}
