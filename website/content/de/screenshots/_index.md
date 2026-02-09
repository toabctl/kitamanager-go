---
title: Bildschirmfotos
weight: 3
---

Ein Rundgang durch die KitaManager Go Benutzeroberflaeche mit den wichtigsten Bildschirmen fuer den taeglichen Einsatz.

---

## Anmeldung

Der Anmeldebildschirm ist der Einstiegspunkt in KitaManager. Benutzer authentifizieren sich mit ihrer E-Mail-Adresse und ihrem Passwort. Nach erfolgreicher Anmeldung wird ein JWT-Token ausgestellt, der die Sitzung aktiv haelt.

{{< figure src="/images/screenshots/login.png" alt="Anmeldeseite" caption="Die Anmeldeseite — geben Sie E-Mail und Passwort ein, um auf das System zuzugreifen." >}}

---

## Dashboard

Nach der Anmeldung bietet das Dashboard einen Ueberblick ueber Ihre Kita auf einen Blick. Die obere Zeile zeigt Zusammenfassungskarten fuer die Gesamtzahl der Organisationen, Mitarbeiter, Kinder und Benutzer. Darunter zeigen die **Schnellstatistiken** Details fuer die aktuell ausgewaehlte Organisation.

Die linke Seitenleiste ist Ihre Hauptnavigation — sie ist unterteilt in systemweite Eintraege (Dashboard, Organisationen, Landesfoerderungen) und organisationsbezogene Eintraege (Benutzer, Gruppen, Mitarbeiter, Kinder, Statistiken, Verguetungsplaene). Ueber das Dropdown in der Seitenleiste koennen Sie zwischen Organisationen wechseln.

{{< figure src="/images/screenshots/dashboard.png" alt="Dashboard" caption="Das Dashboard mit Zusammenfassungskarten und organisationsbezogenen Schnellstatistiken." >}}

---

## Organisationen

Die Organisationsseite listet alle Kita-Einrichtungen auf, auf die Sie Zugriff haben. Jede Zeile zeigt den Organisationsnamen, das Bundesland und ob die Organisation aktuell aktiv ist. Von hier aus koennen Sie neue Organisationen erstellen oder bestehende bearbeiten.

Wenn Sie mehrere Kitas verwalten, erhalten Sie hier das Gesamtbild ueber alle Ihre Einrichtungen.

{{< figure src="/images/screenshots/organizations.png" alt="Organisationsliste" caption="Organisationsuebersicht — verwalten Sie mehrere Kita-Einrichtungen an einem Ort." >}}

---

## Mitarbeiter

Die Mitarbeiterliste zeigt alle Beschaeftigten der ausgewaehlten Organisation. Sie sehen auf einen Blick Name, Geschlecht, Geburtsdatum, Alter, aktuelle Position, Entgeltgruppe und -stufe sowie Wochenstunden jedes Mitarbeiters. Die Aktionsschaltflaechen rechts ermoeglichen es, die Vertragshistorie einzusehen, Details anzuzeigen, den Mitarbeiterdatensatz zu bearbeiten oder einen Mitarbeiter zu entfernen.

Neue Mitarbeiter koennen ueber die Schaltflaeche **+ New Employee** oben rechts hinzugefuegt werden.

{{< figure src="/images/screenshots/employees.png" alt="Mitarbeiterliste" caption="Mitarbeiteruebersicht mit persoenlichen Daten, Position, Entgeltgruppe und Wochenstunden." >}}

---

## Kinder

Die Kinderliste zeigt jedes angemeldete Kind in der ausgewaehlten Organisation. Jede Zeile zeigt Name, Geschlecht, Geburtsdatum, Alter, aktuellen Vertragsstatus, Betreuungseigenschaften (wie halbtag, ganztag, ndh oder integration) und den **automatisch berechneten monatlichen Foerderbetrag** basierend auf der aktiven Landesfoerderungs-Konfiguration.

Dies ist der Bildschirm, den Kita-Leitungen am haeufigsten nutzen — er bietet auf einen Blick ein vollstaendiges Bild von Anmeldungen und Foerderung.

{{< figure src="/images/screenshots/children.png" alt="Kinderliste" caption="Kinderuebersicht mit Anmeldestatus, Betreuungseigenschaften und berechneten Foerderbetraegen." >}}

---

## Landesfoerderung

Die Landesfoerderungsseite ermoeglicht es Administratoren, die landesspezifischen Foerderungsregeln zu konfigurieren, die die automatische Foerderberechnung steuern. Jeder Eintrag repraesentiert eine Foerderungs-Konfiguration fuer ein bestimmtes Bundesland (z.B. "Berlin Kita-Foerderung"). Innerhalb jeder Konfiguration koennen Sie Zeitraeume und eigenschaftsbasierte Foerderbetraege definieren.

Wenn die Vertragseigenschaften eines Kindes mit einem Foerdereintrag uebereinstimmen, wird der entsprechende Monatsbetrag automatisch in der Kinderliste angezeigt.

{{< figure src="/images/screenshots/government-fundings.png" alt="Landesfoerderung" caption="Landesfoerderungs-Konfigurationen — definieren Sie landesspezifische Regeln fuer automatische Foerderberechnungen." >}}
