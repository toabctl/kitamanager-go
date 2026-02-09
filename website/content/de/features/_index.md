---
title: Funktionen
weight: 2
---

KitaManager Go ist eine webbasierte Verwaltungsplattform fuer Kindertagesstaetten (Kitas) in Deutschland. Sie unterstuetzt Einrichtungsleitungen bei den taeglichen Verwaltungsaufgaben — von der Erfassung angemeldeter Kinder mit ihren Vertraegen bis hin zur Personalverwaltung und automatischen Berechnung der Landesfoerderung.

## So funktioniert es

Ein typischer Arbeitsablauf in KitaManager sieht so aus:

```mermaid
flowchart LR
    A["Organisation\nanlegen"] --> B["Mitarbeiter\nerfassen"]
    A --> C["Kinder\nanmelden"]
    B --> D["Vertraege\nzuweisen"]
    C --> E["Betreuungsvertraege\nerstellen"]
    E --> F["Automatische\nFoerderberechnung"]
    style F fill:#22c55e,color:#fff
```

1. **Organisation einrichten** — Ihre Kita mit Name und Bundesland registrieren.
2. **Personal erfassen** — Mitarbeiterdaten eingeben, Positionen zuweisen und Arbeitsvertraege erstellen.
3. **Kinder anmelden** — Kinder mit persoenlichen Daten registrieren und Betreuungsvertraege anlegen.
4. **Foerderung wird automatisch berechnet** — basierend auf den Vertragseigenschaften des Kindes (Betreuungsart, Stunden, besonderer Foerderbedarf) und den Landesfoerderungsregeln.

---

## Organisationsverwaltung

Jede Kita wird im System als **Organisation** abgebildet. Wenn Sie mehrere Einrichtungen betreiben, erhaelt jede eine eigene Organisation mit vollstaendig getrennten Daten.

| Funktion | Beschreibung |
|---|---|
| Mehrere Einrichtungen | Betreiben Sie mehrere Kitas aus einer einzigen KitaManager-Instanz |
| Datentrennung | Kinder, Mitarbeiter und Vertraege sind ihrer Organisation zugeordnet |
| Bundesland-Konfiguration | Jeder Organisation wird ein Bundesland zugewiesen, das die geltenden Foerderregeln bestimmt |
| Gruppen und Bereiche | Organisieren Sie Kinder und Personal in Gruppen innerhalb der Einrichtung |

Administratoren sehen alle ihre Organisationen auf einer Uebersichtsseite und koennen ueber die Seitenleiste zwischen ihnen wechseln.

---

## Personalverwaltung

Das Personalmodul ermoeglicht die Pflege einer vollstaendigen Mitarbeiterdatenbank fuer jede Kita.

### Was Sie pro Mitarbeiter erfassen koennen

| Feld | Beispiel |
|---|---|
| Name, Geschlecht, Geburtsdatum | Anna Mueller, Weiblich, 06.05.2000 |
| Position | Erzieher, Kinderpfleger, Gruppenleitung |
| Entgeltgruppe und Stufe | S8a / Stufe 3 |
| Wochenstunden | 39 Stunden |
| Vertragszeitraum | 01.01.2024 — 31.12.2025 |

### Arbeitsvertraege

Jeder Mitarbeiter kann im Laufe der Zeit einen oder mehrere **Arbeitsvertraege** haben. Vertraege definieren Position, Entgeltgruppe, Wochenstunden und Gueltigkeitszeitraum. Das System stellt sicher, dass sich Vertraege fuer denselben Mitarbeiter nicht ueberschneiden.

### Verguetungsplaene

Verguetungsplaene definieren die in Ihrer Einrichtung verwendeten Entgeltgruppen und Stufen (z.B. die TVoeD-SuE-Tabelle, die in deutschen oeffentlichen Kindertagesstaetten ueblich ist). Wenn Sie einem Mitarbeitervertrag eine Gruppe und Stufe zuweisen, verfolgt das System die Entwicklung.

---

## Kinderverwaltung

Das Kindermodul verfolgt jedes in Ihrer Kita angemeldete Kind sowie seine Betreuungsvertraege und Foerderung.

### Was Sie pro Kind erfassen koennen

| Feld | Beispiel |
|---|---|
| Name, Geschlecht, Geburtsdatum | Laura Lange, Weiblich, 27.03.2025 |
| Aktueller Vertragsstatus | Aktiv, Bevorstehend oder Beendet |
| Betreuungseigenschaften | halbtag, ganztag, teilzeit, integration, ndh |
| Berechnete monatliche Foerderung | 1.215,45 EUR |

### Betreuungsvertraege

Jedes Kind hat einen oder mehrere **Betreuungsvertraege**, die den Anmeldezeitraum und die Art der Betreuung festlegen. Vertragseigenschaften sind Merkmale, die die Betreuungsvereinbarung beschreiben:

| Eigenschaft | Bedeutung |
|---|---|
| `halbtag` | Halbtagsbetreuung |
| `ganztag` | Ganztagsbetreuung |
| `teilzeit` | Teilzeitbetreuung |
| `ndh` | Nichtdeutsche Herkunftssprache |
| `integration a/b` | Integrationsstufen |

Diese Eigenschaften bestimmen direkt, wie viel Landesfoerderung die Kita fuer jedes Kind erhaelt (siehe [Landesfoerderung](#landesfoerderung) unten).

---

## Landesfoerderung

Eine der Kernfunktionen von KitaManager ist die automatische Berechnung der staatlichen Kita-Foerderung basierend auf den Regeln des jeweiligen Bundeslandes.

### So funktioniert die Foerderberechnung

```mermaid
flowchart TD
    A["Betreuungsvertrag"] --> B["Vertragseigenschaften\n(z.B. ganztag + ndh)"]
    B --> C["Abgleich mit\nFoerderregeln"]
    D["Landesfoerderungs-Konfiguration\n(z.B. Berlin 2024)"] --> C
    C --> E["Monatlicher Foerderbetrag\n(z.B. 1.633,64 EUR)"]
    style E fill:#22c55e,color:#fff
```

1. Jeder Betreuungsvertrag hat **Eigenschaften**, die die Betreuungsart beschreiben.
2. Das System sucht den passenden **Foerdereintrag** aus den konfigurierten Landesfoerderungsregeln.
3. Der resultierende **Monatsbetrag** wird direkt in der Kinderliste angezeigt.

### Foerderungs-Konfiguration

Die Foerderung wird pro Bundesland und Zeitraum konfiguriert. Jeder Foerdereintrag ordnet eine Kombination von Eigenschaften einem monatlichen Betrag in Euro zu:

| Eigenschaften | Monatsbetrag |
|---|---|
| halbtag | 1.215,45 EUR |
| halbtag + ndh | 1.318,11 EUR |
| ganztag | 1.909,61 EUR |
| ganztag + integration a | 3.566,41 EUR |
| teilzeit + ndh | 1.633,64 EUR |

Foerderzeitraeume koennen aktualisiert werden, wenn sich die Landessaetze aendern, ohne historische Daten zu beeinflussen.

---

## Benutzerrollen und Zugriffskontrolle

KitaManager verwendet ein rollenbasiertes Zugriffskontrollsystem (RBAC), das sicherstellt, dass Benutzer nur auf die fuer ihre Rolle und Organisation relevanten Daten zugreifen koennen.

### Rollenuebersicht

| Rolle | Geltungsbereich | Mitarbeiter verwalten | Kinder verwalten | Foerderung verwalten | Benutzer verwalten |
|---|---|---|---|---|---|
| **Superadmin** | Alle Organisationen | Ja | Ja | Ja | Ja |
| **Admin** | Zugewiesene Org(s) | Ja | Ja | Ja | Ja |
| **Manager** | Zugewiesene Org(s) | Ja | Ja | Nein | Nein |
| **Mitglied** | Zugewiesene Org(s) | Nur lesen | Nur lesen | Nein | Nein |

- **Superadmin** ist der systemweite Administrator, der alle Organisationen und Benutzer verwalten kann.
- **Admin** hat volle Kontrolle innerhalb einer oder mehrerer zugewiesener Organisationen.
- **Manager** kuemmert sich um das Tagesgeschaeft wie Mitarbeiter- und Kinderverwaltung.
- **Mitglied** kann Daten einsehen, aber keine Aenderungen vornehmen.

Alle Datenaenderungen werden in einem Audit-Log fuer Compliance-Zwecke protokolliert.

---

## Dashboard und Berichte

Nach der Anmeldung sehen Benutzer ein **Dashboard**, das einen schnellen Ueberblick ueber ihre Kita bietet:

- Gesamtzahl der Organisationen, Mitarbeiter, Kinder und Benutzer
- Schnellstatistiken bezogen auf die aktuell ausgewaehlte Organisation
- Ein-Klick-Navigation zu allen Verwaltungsbereichen ueber die Seitenleiste

Die Seitenleiste bietet direkten Zugriff auf:

| Menuepunkt | Zweck |
|---|---|
| Dashboard | Ueberblick und Kennzahlen |
| Organisationen | Kita-Einrichtungen verwalten |
| Landesfoerderungen | Landesfoerderungsregeln konfigurieren |
| Benutzer | Benutzerkonten verwalten |
| Gruppen | Organisationsgruppen verwalten |
| Mitarbeiter | Personaldatenbank und Vertraege |
| Kinder | Anmeldungen und Betreuungsvertraege |
| Statistiken | Berichte und Datenanalyse |
| Verguetungsplaene | Entgeltgruppen-Definitionen |
